package api

type DatabaseTester interface {
	CreateDB(dbName string)
	CreateTable(dbName, tableName string)
	InsertTable(dbName, tableName, data string)
	QueryAll(dbName, tableName string)
}
