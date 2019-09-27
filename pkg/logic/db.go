package logic

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/shiywang/tist/pkg/base/mysql"
	"github.com/shiywang/tist/pkg/util"
)

type MysqlDB struct {
	sql *mysql.SQLConnPool
}

const (
	tableFieldCount    = 3
	localhost          = "127.0.0.1"
	defaultPort        = "4000"
	defaultUser        = "root"
	defaultPassWord    = ""
	defaultCharSet     = "utf8"
	defaultMaxOpenConn = 20
	defaultMaxIdleConn = 20
)

func (m *MysqlDB) CreateDB(dbName string) {
	m.sql = mysql.CreateMySQLClient(localhost, defaultPort, dbName, defaultUser, defaultPassWord, defaultCharSet,
		defaultMaxOpenConn, defaultMaxIdleConn)

	if m.sql != nil {
		return
	}
	util.CheckErr(errors.New("CreateMySQLClient return nil"))
}

func (m *MysqlDB) CreateTable(dbName, tableName string) {
	if m.sql != nil {
		if _, err := m.sql.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)); err != nil {
			util.CheckErr(err)
		}

		buf := new(bytes.Buffer)
		s := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (TEST_KEY INT AUTO_INCREMENT PRIMARY KEY", tableName)
		buf.WriteString(s)

		for i := 0; i < tableFieldCount; i++ {
			buf.WriteString(fmt.Sprintf(", FIELD%d VARCHAR(%d)", i, 100))
		}

		buf.WriteString(");")

		if _, err := m.sql.Exec(buf.String()); err != nil {
			util.CheckErr(err)
		}
		return
	}
	util.CheckErr(errors.New("CreateMySQLClient return nil"))
}

func (m *MysqlDB) InsertTable(dbName, tableName, data string) {
	if m.sql != nil {
		buf := new(bytes.Buffer)

		buf.WriteString("INSERT INTO ")
		buf.WriteString(tableName)
		buf.WriteString(" (FIELD0")

		for i := 1; i < tableFieldCount; i++ {
			buf.WriteString(fmt.Sprintf(", FIELD%d", i))
		}
		buf.WriteString(") VALUES (")

		for i := 0; i < tableFieldCount-1; i++ {
			if i == 0 {
				buf.WriteString(fmt.Sprintf("\"%s\"", data))
			}
			buf.WriteString(fmt.Sprintf(" ,\"%s\"", data))
		}

		buf.WriteByte(')')

		if _, err := m.sql.Insert(buf.String()); err != nil {
			util.CheckErr(err)
		}
		return
	}
	util.CheckErr(errors.New("CreateMySQLClient return nil"))
}

func (m *MysqlDB) QueryAll(dbName, tableName string) {
	if m.sql != nil {
		var out []map[string]interface{}
		var err error
		sql := "select * from " + tableName
		if out, err = m.sql.Query(sql); err != nil {
			util.CheckErr(err)
		}
		fmt.Println("return: ", out)
		return
	}
	util.CheckErr(errors.New("CreateMySQLClient return nil"))

}
