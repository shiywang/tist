/*

This code is a reference from https://github.com/hopehook/golang-db
with slightly modification

*/

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/color"
	"github.com/shiywang/tist/pkg/util"
)

// SQLConnPool is DB pool struct
type SQLConnPool struct {
	DriverName     string
	DataSourceName string
	MaxOpenConns   int
	MaxIdleConns   int
	sql            *sql.DB
}

// CreateMySQLClient func init DB pool
func CreateMySQLClient(host, port, database, user, password, charset string, maxOpenConns, maxIdleConns int) *SQLConnPool {
	var db *SQLConnPool
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&autocommit=true", user, password, host, port, database, charset)
	db = &SQLConnPool{
		DriverName:     "mysql",
		DataSourceName: dataSourceName,
		MaxOpenConns:   maxOpenConns,
		MaxIdleConns:   maxIdleConns,
	}
	if err := db.Open(); err != nil {
		util.CheckErr(err)
	}
	return db
}

func (p *SQLConnPool) RetryPing() error {
	yellow := color.FgYellow.Render
	var err error
	wait := 5 * time.Second
	for attempts := 0; attempts < 5; attempts++ {
		fmt.Println(yellow(fmt.Sprintf("retry ping for %d time(s).....", attempts+1)))
		if err = p.sql.Ping(); err == nil {
			return nil
		}
		wait = wait * 2
		time.Sleep(wait)
	}
	return err
}

func (p *SQLConnPool) Open() error {
	var err error
	p.sql, err = sql.Open(p.DriverName, p.DataSourceName)
	if err != nil {
		return err
	}
	if err = p.RetryPing(); err != nil {
		return err
	}
	p.sql.SetMaxOpenConns(p.MaxOpenConns)
	p.sql.SetMaxIdleConns(p.MaxIdleConns)
	return nil
}

// Close pool
func (p *SQLConnPool) Close() error {
	return p.sql.Close()
}

// Get via pool
func (p *SQLConnPool) Get(queryStr string, args ...interface{}) (map[string]interface{}, error) {
	results, err := p.Query(queryStr, args...)
	if err != nil {
		return map[string]interface{}{}, err
	}
	if len(results) <= 0 {
		return map[string]interface{}{}, sql.ErrNoRows
	}
	if len(results) > 1 {
		return map[string]interface{}{}, errors.New("sql: more than one rows")
	}
	return results[0], nil
}

// Query via pool
func (p *SQLConnPool) Query(queryStr string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := p.sql.Query(queryStr, args...)
	if err != nil {
		log.Println(err)
		return []map[string]interface{}{}, err
	}
	defer rows.Close()
	columns, err := rows.ColumnTypes()
	scanArgs := make([]interface{}, len(columns))
	values := make([]sql.RawBytes, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	rowsMap := make([]map[string]interface{}, 0, 10)
	for rows.Next() {
		rows.Scan(scanArgs...)
		rowMap := make(map[string]interface{})
		for i, value := range values {
			rowMap[columns[i].Name()] = bytes2RealType(value, columns[i])
		}
		rowsMap = append(rowsMap, rowMap)
	}
	if err = rows.Err(); err != nil {
		return []map[string]interface{}{}, err
	}
	return rowsMap, nil
}

func (p *SQLConnPool) Exec(sqlStr string, args ...interface{}) (sql.Result, error) {
	return p.sql.Exec(sqlStr, args...)
}

// Update via pool
func (p *SQLConnPool) Update(updateStr string, args ...interface{}) (int64, error) {
	result, err := p.Exec(updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

// Insert via pool
func (p *SQLConnPool) Insert(insertStr string, args ...interface{}) (int64, error) {
	result, err := p.Exec(insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	return lastId, err

}

// Delete via pool
func (p *SQLConnPool) Delete(deleteStr string, args ...interface{}) (int64, error) {
	result, err := p.Exec(deleteStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

// bytes2RealType is to convert db type to code type
func bytes2RealType(src []byte, column *sql.ColumnType) interface{} {
	srcStr := string(src)
	var result interface{}
	switch column.DatabaseTypeName() {
	case "BIT", "TINYINT", "SMALLINT", "INT":
		result, _ = strconv.ParseInt(srcStr, 10, 64)
	case "BIGINT":
		result, _ = strconv.ParseUint(srcStr, 10, 64)
	case "CHAR", "VARCHAR",
		"TINY TEXT", "TEXT", "MEDIUM TEXT", "LONG TEXT",
		"TINY BLOB", "MEDIUM BLOB", "BLOB", "LONG BLOB",
		"JSON", "ENUM", "SET",
		"YEAR", "DATE", "TIME", "TIMESTAMP", "DATETIME":
		result = srcStr
	case "FLOAT", "DOUBLE", "DECIMAL":
		result, _ = strconv.ParseFloat(srcStr, 64)
	default:
		result = nil
	}
	return result
}
