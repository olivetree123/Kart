package database

import (
	"Kart/utils"
	"fmt"
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

type Condition struct {
	Field    string
	Value    string
	Operator string
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

func (conn *Connection) ParseCondition(tableID [32]byte, condition string) []Condition {
	var conditions []Condition
	conds := strings.Split(condition, "and")
	for _, cond := range conds {
		cond = strings.TrimSpace(cond)
		operator := FindOperator(cond)
		ds := strings.Split(cond, operator)
		if len(ds) != 2 {
			panic("Invalid condition.")
		}
		column := conn.MetaDataBase.FindColumnByName(tableID, ds[0])
		if column == nil {
			panic("Invalid column.")
		}
		value := ds[1]
		if utils.SliceToString(column.Type[:]) == "bool" {
			fmt.Println("column type is bool")
			if value == "true" {
				value = "1"
			} else {
				value = "0"
			}
		}
		c := Condition{
			Field:    ds[0],
			Value:    value,
			Operator: operator,
		}
		conditions = append(conditions, c)
	}
	return conditions
}

func (conn *Connection) SelectOne(tableName string, condition string) map[string]string {
	table := conn.MetaDataBase.FindTableByName(conn.DBID, tableName)
	if table == nil {
		panic("Table does not exist.")
	}
	fmt.Println("tableName = ", tableName, utils.SliceToString(table.Name[:]))
	conditions := conn.ParseCondition(table.ID, condition)
	return conn.DataDataBase.SelectOneData(table.ID, conditions)
}

func (conn *Connection) Select(tableName string, condition string) []map[string]string {
	table := conn.MetaDataBase.FindTableByName(conn.DBID, tableName)
	if table == nil {
		panic("Table does not exist.")
	}
	conditions := conn.ParseCondition(table.ID, condition)
	return conn.DataDataBase.SelectData(table.ID, conditions)
}

func (conn *Connection) Update(tableName string, condition string, data map[string]string) {
	table := conn.MetaDataBase.FindTableByName(conn.DBID, tableName)
	if table == nil {
		panic("Table does not exist.")
	}
	conditions := conn.ParseCondition(table.ID, condition)
	conn.DataDataBase.UpdateData(table.ID, conditions, data)
}

func (conn *Connection) Delete(tableName string, condition string) {
	table := conn.MetaDataBase.FindTableByName(conn.DBID, tableName)
	if table == nil {
		panic("Table does not exist.")
	}
	conditions := conn.ParseCondition(table.ID, condition)
	conn.DataDataBase.DeleteData(table.ID, conditions)
}
