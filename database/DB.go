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
	Length uint8
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
	Data    map[string]string
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

func NewDataDB(metaDB *MetaDB) *DataDB {
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
	dataDB.LoadDataDB(metaDB)
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

func (db *DataDB) LoadDataDB(metaDB *MetaDB) {
	println("load data db")
	size, err := db.FileHandler.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	if size == 0 {
		return
	}
	fmt.Println("size = ", size)
	var offset int64
	for size > 0 {
		// 1. 读 TableID
		tableID := make([]byte, 32)
		_, err := db.FileHandler.ReadAt(tableID, offset)
		if err != nil {
			panic(err)
		}
		size -= 32
		offset += 32
		// 2. 读 DataID
		dataID := make([]byte, 32)
		_, err = db.FileHandler.ReadAt(dataID, offset)
		if err != nil {
			panic(err)
		}
		size -= 32
		offset += 32
		// 3. 读 Status
		status := make([]byte, 1)
		_, err = db.FileHandler.ReadAt(status, offset)
		if err != nil {
			panic(err)
		}
		size -= 1
		offset += 1
		// 4. 读 ColumnID
		columnID := make([]byte, 32)
		_, err = db.FileHandler.ReadAt(columnID, offset)
		if err != nil {
			panic(err)
		}
		size -= 32
		offset += 32
		// 5. 读 ColumnLength
		columnLenBytes := make([]byte, 1)
		_, err = db.FileHandler.ReadAt(columnLenBytes, offset)
		if err != nil {
			panic(err)
		}
		size -= 1
		offset += 1
		//columnLen := binary.LittleEndian.Uint16(columnLenBytes)
		var buf bytes.Buffer
		_, err = buf.Write(columnLenBytes)
		if err != nil {
			panic(err)
		}
		columnLen, _ := binary.ReadUvarint(&buf)
		// 6. 读 ColumnData
		colData := make([]byte, columnLen)
		_, err = db.FileHandler.ReadAt(colData, offset)
		if err != nil {
			panic(err)
		}
		size -= int64(columnLen)
		offset += int64(columnLen)

		found := false
		for _, td := range db.TbData {
			if td.TableID == utils.SliceToUUID(tableID) && td.Data["ID"] == utils.SliceToString(dataID) {
				found = true
				td.Data[utils.SliceToString(columnID)] = utils.SliceToString(colData)
			}
		}
		if !found {
			data := make(map[string]string)
			data["ID"] = utils.SliceToString(dataID)
			column := metaDB.FindColumnByID(utils.SliceToUUID(tableID), utils.SliceToUUID(columnID))
			data[utils.SliceToString(column.Name[:])] = utils.SliceToString(colData)
			tbData := &TableData{
				TableID: utils.SliceToUUID(tableID),
				Data:    data,
			}
			db.TbData = append(db.TbData, tbData)
		}
		//tbData := &TableData{
		//	TableID: utils.SliceToUUID(tableID),
		//	Data:    data,
		//}
		//db.TbData = append(db.TbData, tbData)
	}
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
		db.Columns = append(db.Columns, data.(*ColumnMetaData))
	}
}

func (db *DataDB) AddData(tableID [32]byte, dataID [32]byte, columns []*ColumnMetaData, model Model) {
	// 数据文件结构：TableID + DataID + Status + ColumnID + ColumnLength + ColumnData
	// 各个字段字节数分别为：32 + 32 + 1 + 32 + 8 + [ColumnLen]
	var buf bytes.Buffer
	for _, column := range columns {
		err := binary.Write(&buf, binary.LittleEndian, tableID)
		if err != nil {
			panic(err)
		}
		err = binary.Write(&buf, binary.LittleEndian, dataID)
		if err != nil {
			panic(err)
		}
		var status = [1]byte{1}
		err = binary.Write(&buf, binary.LittleEndian, status)
		if err != nil {
			panic(err)
		}
		err = binary.Write(&buf, binary.LittleEndian, column.ID)
		if err != nil {
			panic(err)
		}
		fmt.Println("Write column.Length = ", column.Length)
		err = binary.Write(&buf, binary.LittleEndian, column.Length)
		if err != nil {
			panic(err)
		}
		for _, field := range model.Fields() {
			if field.GetName() != utils.SliceToString(column.Name[:]) {
				continue
			}
			err = binary.Write(&buf, binary.LittleEndian, field.Bytes())
			if err != nil {
				panic(err)
			}
		}
	}
	err := binary.Write(db.FileHandler, binary.LittleEndian, buf.Bytes())
	if err != nil {
		panic(err)
	}
	d := &TableData{
		TableID: tableID,
		Data:    ModelToMap(model),
	}
	db.TbData = append(db.TbData, d)
}

func (db *DataDB) SelectData(tableID [32]byte, columns []*ColumnMetaData, queryMap map[string]interface{}) []map[string]string {
	var result []map[string]string
	fmt.Println("select tableID = ", tableID)
	for _, d := range db.TbData {
		fmt.Println("d.TableID = ", d.TableID)
		if d.TableID != tableID {
			continue
		}
		fmt.Println("dataMap = ", d.Data)
		fmt.Println("queryMap = ", queryMap)
		flag := true
		for key, value := range queryMap {
			fmt.Println("key = ", key, "value = ", d.Data[key])
			if d.Data[key] != value {
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

func (db *MetaDB) FindColumnByName(tableID [32]byte, columnName string) *ColumnMetaData {
	for _, column := range db.Columns {
		if columnName == utils.SliceToString(column.Name[:]) {
			return column
		}
	}
	return nil
}

func (db *MetaDB) FindColumnByID(tableID [32]byte, columnID [32]byte) *ColumnMetaData {
	for _, column := range db.Columns {
		if column.TbID == tableID && column.ID == columnID {
			return column
		}
	}
	return nil
}

func (db *MetaDB) AddColumn(tableID [32]byte, columnName string, tp string, length int) {
	for _, column := range db.Columns {
		if utils.SliceToString(column.Name[:]) == columnName {
			fmt.Println("column already exists: ", columnName)
			return
		}
	}
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
		Length: uint8(length),
	}
	db.AddData(utils.UUIDToString(uid), MagicCOLUMN, m)
}
