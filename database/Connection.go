package database

import (
	"fmt"
	"kart/utils"
	"reflect"
	"strings"
	//"strings"
)

type Connection struct {
	DBID         string
	DBName       string
	MetaDataBase *MetaDB
	DataDataBase *DataDB
}

func NewConnection(dbName string) *Connection {
	fmt.Println("111")
	metaDB := NewMetaDB()
	fmt.Println("begin to create db")
	db := metaDB.CreateDB(dbName)
	conn := &Connection{
		DBID:         utils.UUIDToString(db.ID),
		DBName:       dbName,
		MetaDataBase: metaDB,
		DataDataBase: NewDataDB(metaDB),
	}
	return conn
}

func (conn *Connection) CreateTable(model Model) {
	tableName := reflect.TypeOf(model).Name()
	table := conn.MetaDataBase.CreateTable(conn.DBID, tableName)
	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		//fmt.Println(t.Field(i).Name, t.Field(i).Type.Name(), t.Field(i).Tag.Get("orm"))
		length := 0
		switch t.Field(i).Type.Name() {
		case "StringField":
			length = utils.GetLenFromTag(t.Field(i).Tag.Get("orm"))
			break
		case "UUIDField":
			length = 32
			break
		case "BooleanField":
			length = 1
			break
		case "IntegerField":
			length = 20
			break
		}
		conn.MetaDataBase.AddColumn(table.ID, t.Field(i).Name, "string", length)
	}
}

func (conn *Connection) Insert(tableName string, data Model) {
	// 需要验证 data 与数据库中的表是否结构相同
	//dataMap := utils.StructPtr2Map(data)
	table := conn.MetaDataBase.FindTableByName(conn.DBID, tableName)
	if table == nil {
		panic("Table does not exist.")
	}
	columns := conn.MetaDataBase.FindColumnByTable(table.ID)
	colMap := make(map[string]interface{})
	for _, column := range columns {
		d := make(map[string]interface{})
		d["Type"] = utils.SliceToString(column.Type[:])
		d["Length"] = column.Length
		colMap[utils.SliceToString(column.Name[:])] = d
	}
	for _, field := range data.Fields() {
		fieldName := field.(Field).GetName()
		if _, found := colMap[fieldName]; !found {
			info := fmt.Sprintf("key = %s not found in table %s", fieldName, tableName)
			panic(info)
		}
	}
	conn.DataDataBase.AddData(table.ID, data.GetID(), columns, data)
}

func (conn *Connection) Select(tableName string, condition string) []map[string]string {
	table := conn.MetaDataBase.FindTableByName(conn.DBID, tableName)
	if table == nil {
		panic("Table does not exist.")
	}
	columns := conn.MetaDataBase.FindColumnByTable(table.ID)
	queryMap := make(map[string]interface{})
	conds := strings.Split(condition, "and")
	for _, cond := range conds {
		cond = strings.TrimSpace(cond)
		ds := strings.Split(cond, "=")
		if len(ds) != 2 {
			panic("Invalid condition.")
		}
		column := conn.MetaDataBase.FindColumnByName(table.ID, ds[0])
		if column == nil {
			panic("Invalid column.")
		}
		queryMap[utils.SliceToString(column.ID[:])] = ds[1]
	}
	return conn.DataDataBase.SelectData(table.ID, columns, queryMap)
}
