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
	metaDB := NewMetaDB()
	db := metaDB.CreateDB(dbName)
	conn := &Connection{
		DBID:         utils.UUIDToString(db.ID),
		DBName:       dbName,
		MetaDataBase: metaDB,
		DataDataBase: NewDataDB(metaDB),
	}
	return conn
}

func (conn *Connection) CreateTable(model interface{}) {
	tableName := reflect.TypeOf(model).Name()
	table := conn.MetaDataBase.CreateTable(conn.DBID, tableName)
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)
	for i := 0; i < t.NumField(); i++ {
		length := v.Field(i).Interface().(Field).GetLength()
		if t.Field(i).Type.Name() == "StringField" {
			length2 := utils.GetLenFromTag(t.Field(i).Tag.Get("kart"))
			if length2 > 0 {
				length = length2
			}
		}
		conn.MetaDataBase.AddColumn(table.ID, t.Field(i).Name, v.Field(i).Interface().(Field).GetType(), length)
	}
}

func (conn *Connection) Insert(tableName string, data interface{}) {
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
	for _, field := range GetModelFields(data) {
		fieldName := field.(Field).GetName()
		if _, found := colMap[fieldName]; !found {
			info := fmt.Sprintf("key = %s not found in table %s", fieldName, tableName)
			panic(info)
		}
	}
	dataID, err := GetModelID(data)
	if err != nil {
		panic(err)
	}
	conn.DataDataBase.AddData(table.ID, utils.StringToUUID(dataID), columns, data)
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
		value := ds[1]
		fmt.Println("column type = ", utils.SliceToString(column.Type[:]))
		if utils.SliceToString(column.Type[:]) == "bool" {
			fmt.Println("column type is bool")
			if value == "true" {
				value = "1"
			} else {
				value = "0"
			}
		}
		//queryMap[utils.SliceToString(column.ID[:])] = value
		queryMap[ds[0]] = value
	}
	return conn.DataDataBase.SelectData(table.ID, columns, queryMap)
}
