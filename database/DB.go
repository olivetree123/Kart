package database

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"kart/config"
	"kart/utils"
	"os"
	"path/filepath"
	"unsafe"
)

const MagicDB int8 = 1
const MagicTABLE int8 = 2
const MagicCOLUMN int8 = 3

type DBMetaData struct {
	ID   [32]byte
	Name [50]byte
}

type TableMetaData struct {
	ID   [32]byte
	Name [50]byte
	// DbID 所属 DB 的 ID
	DbID [32]byte
}

type ColumnMetaData struct {
	ID   [32]byte
	Name [50]byte
	// TbID 所属 Table 的 ID
	TbID [32]byte
	// Type 字段类型，有效值有：string/int/bool
	Type [50]byte
	// Length 字段长度，只有在 Type == string 时有效
	Length int32
}

// MetaDB 存储所有DB元数据的数据库
type MetaDB struct {
	FilePath    string
	FileHandler *os.File
	DBs         []*DBMetaData
	Tables      []*TableMetaData
	Columns     []*ColumnMetaData
}

type DataDB struct {
	FilePath    string
	FileHandler *os.File
	TbData      []*TableData
}

type TableData struct {
	TableID [32]byte
	Data    []byte
}

func GetObjectByMagicNumber(magicNum int) (interface{}, int) {
	if magicNum == int(MagicDB) {
		return &DBMetaData{}, binary.Size(DBMetaData{})
	} else if magicNum == int(MagicTABLE) {
		return &TableMetaData{}, binary.Size(TableMetaData{})
	} else if magicNum == int(MagicCOLUMN) {
		return &ColumnMetaData{}, binary.Size(ColumnMetaData{})
	}
	return nil, 0
}

func NewMetaDB() *MetaDB {
	metaFilePath := filepath.Join(config.Config.GetString("FilePath"), config.Config.GetString("MetaFileName"))
	f, err := os.OpenFile(metaFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	metaDB := &MetaDB{
		FilePath:    metaFilePath,
		FileHandler: f,
		DBs:         nil,
		Tables:      nil,
		Columns:     nil,
	}
	metaDB.LoadMetaDB()
	return metaDB
}

func NewDataDB() *DataDB {
	metaFilePath := filepath.Join(config.Config.GetString("FilePath"), config.Config.GetString("DataFileName"))
	f, err := os.OpenFile(metaFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	var data []*TableData
	dataDB := &DataDB{
		FilePath:    metaFilePath,
		FileHandler: f,
		TbData:      data,
	}
	dataDB.LoadDataDB()
	return dataDB
}

func (db *MetaDB) LoadMetaDB() {
	size, err := db.FileHandler.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	if size == 0 {
		return
	}
	fmt.Println("size = ", size)
	//length := unsafe.Sizeof(MetaDB{})
	var offset int64
	for size > 0 {
		var magicNum int8
		magicBytes := make([]byte, unsafe.Sizeof(int8(1)))
		_, err := db.FileHandler.ReadAt(magicBytes, offset)
		if err != nil {
			panic(err)
		}
		buf := bytes.NewReader(magicBytes)
		err = binary.Read(buf, binary.LittleEndian, &magicNum)
		if err != nil {
			panic(err)
		}
		size -= 1
		offset += 1
		st, length := GetObjectByMagicNumber(int(magicNum))
		indexBytes := make([]byte, length)
		_, err = db.FileHandler.ReadAt(indexBytes, offset)
		if err != nil {
			panic(err)
		}
		buf = bytes.NewReader(indexBytes)
		err = binary.Read(buf, binary.LittleEndian, st)
		if err != nil {
			panic(err)
		}
		if int8(magicNum) == MagicDB {
			db.DBs = append(db.DBs, st.(*DBMetaData))
		} else if int8(magicNum) == MagicTABLE {
			db.Tables = append(db.Tables, st.(*TableMetaData))
		} else if int8(magicNum) == MagicCOLUMN {
			db.Columns = append(db.Columns, st.(*ColumnMetaData))
		}
		size -= int64(length)
		offset += int64(length)
	}
}

func (db *DataDB) LoadDataDB() {
	println("load data db")
}

func (db *MetaDB) AddData(id string, magicNum int8, data interface{}) {
	err := binary.Write(db.FileHandler, binary.LittleEndian, magicNum)
	if err != nil {
		panic(err)
	}
	if int8(magicNum) == MagicDB {
		err = binary.Write(db.FileHandler, binary.LittleEndian, data.(*DBMetaData))
		if err != nil {
			panic(err)
		}
		db.DBs = append(db.DBs, data.(*DBMetaData))
	} else if int8(magicNum) == MagicTABLE {
		err = binary.Write(db.FileHandler, binary.LittleEndian, data.(*TableMetaData))
		if err != nil {
			panic(err)
		}
		db.Tables = append(db.Tables, data.(*TableMetaData))
	} else if int8(magicNum) == MagicCOLUMN {
		err = binary.Write(db.FileHandler, binary.LittleEndian, data.(*ColumnMetaData))
		if err != nil {
			panic(err)
		}
		fmt.Println("add column = ", utils.SliceToString(data.(*ColumnMetaData).Name[:]))
		db.Columns = append(db.Columns, data.(*ColumnMetaData))
	}
}

func (db *DataDB) AddData(tableID [32]byte, data []byte) {
	err := binary.Write(db.FileHandler, binary.LittleEndian, tableID)
	if err != nil {
		panic(err)
	}
	err = binary.Write(db.FileHandler, binary.LittleEndian, data)
	if err != nil {
		panic(err)
	}
	d := &TableData{
		TableID: tableID,
		Data:    data,
	}
	db.TbData = append(db.TbData, d)
}

func (db *DataDB) SelectData(tableID [32]byte, columns []*ColumnMetaData, queryMap map[string]interface{}, colMap map[string]interface{}) []interface{} {
	var result []interface{}
	for _, d := range db.TbData {
		if d.TableID != tableID {
			continue
		}
		offset := 0
		dataMap := make(map[string]string)
		fmt.Println("data length = ", len(d.Data))
		for _, column := range columns {
			fmt.Println("column Name = ", utils.SliceToString(column.Name[:]), "length = ", column.Length)
			fmt.Println("endpoint = ", offset+int(column.Length))
			value := utils.SliceToString(d.Data[offset : offset+int(column.Length)])
			dataMap[utils.SliceToString(column.Name[:])] = value
			offset += int(column.Length)
		}
		fmt.Println("dataMap = ", dataMap)
		fmt.Println("queryMap = ", queryMap)
		//dMap := utils.StructPtr2Map(d.Data)
		flag := true
		for key, value := range queryMap {
			//buf := bytes.NewBuffer(reflect.ValueOf(dMap[key]).Interface().([]byte))
			//fmt.Println(buf.String())
			//fmt.Println("type of key = ", reflect.TypeOf(dMap[key]), reflect.ValueOf(dMap[key]), reflect.ValueOf(utils.StringToSlice2(value.(string), 32)))
			//fmt.Println("compare: ", key, reflect.ValueOf(dMap[key]), reflect.ValueOf(utils.StringToSlice2(value.(string), 32)))
			//if reflect.ValueOf(dMap[key]) != reflect.ValueOf(utils.StringToSlice2(value.(string), 32)) {
			//	flag = false
			//	break
			//}
			//fmt.Println("Yes, Equal !")
			if dataMap[key] != value {
				flag = false
				break
			}
			fmt.Println("Yes, Equal !")
		}
		if flag {
			result = append(result, d.Data)
		}
	}
	return result
}

func (db *MetaDB) CreateDB(name string) *DBMetaData {
	r := db.FindDBByName(name)
	if r != nil {
		return r
	}
	uid := utils.NewUUID()
	var nm [50]byte
	copy(nm[:], name)
	fmt.Println("nm length = ", len(nm))
	m := &DBMetaData{
		ID:   uid,
		Name: nm,
	}
	db.AddData(utils.UUIDToString(uid), MagicDB, m)
	return m
}

func (db *MetaDB) CreateTable(dbID string, tableName string) *TableMetaData {
	r := db.FindTableByName(dbID, tableName)
	if r != nil {
		return r
	}
	uid := utils.NewUUID()
	var tb [50]byte
	copy(tb[:], tableName)
	m := &TableMetaData{
		ID:   uid,
		Name: tb,
		DbID: utils.StringToUUID(dbID),
	}
	db.AddData(utils.UUIDToString(uid), MagicTABLE, m)
	return m
}

func (db *MetaDB) FindDBByName(dbName string) *DBMetaData {
	for _, db := range db.DBs {
		if dbName == utils.SliceToString(db.Name[:]) {
			return db
		}
	}
	return nil
}

func (db *MetaDB) FindTableByName(dbID string, tableName string) *TableMetaData {
	for _, table := range db.Tables {
		if utils.SliceToString(table.Name[:]) == tableName &&
			utils.UUIDToString(table.DbID) == dbID {
			return table
		}
	}
	return nil
}

func (db *MetaDB) FindColumnByTable(tableID [32]byte) []*ColumnMetaData {
	var result []*ColumnMetaData
	for _, column := range db.Columns {
		if column.TbID == tableID {
			result = append(result, column)
		}
	}
	return result
}

func (db *MetaDB) AddColumn(tableID [32]byte, columnName string, tp string, length int) {
	uid := utils.NewUUID()
	var cm [50]byte
	var t [50]byte
	copy(cm[:], columnName)
	copy(t[:], tp)
	m := &ColumnMetaData{
		ID:     uid,
		Name:   cm,
		TbID:   tableID,
		Type:   t,
		Length: int32(length),
	}
	db.AddData(utils.UUIDToString(uid), MagicCOLUMN, m)
}
