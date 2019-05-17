package database

import (
	"Kart/config"
	"Kart/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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

type ColumnData struct {
	PKValue string
	Value   string
	Offset  int64
	Length  int
	Status  bool
}

type TableData struct {
	TableID [32]byte
	Data    map[string]*ColumnData
}

func ColumnDataToMap(data map[string]*ColumnData) map[string]string {
	result := make(map[string]string)
	for key, d := range data {
		result[key] = d.Value
	}
	return result
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
	f, err := os.OpenFile(metaFilePath, os.O_RDWR|os.O_CREATE, 0755)
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
	//for _, table := range db.Tables {
	//	fmt.Println("Table = ", utils.SliceToString(table.Name[:]), utils.SliceToString(table.ID[:]))
	//}
	//for _, column := range db.Columns {
	//	fmt.Println("Column = ", utils.SliceToString(column.TbID[:]), utils.SliceToString(column.ID[:]), utils.SliceToString(column.Name[:]))
	//}
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
		st := true
		if utils.SliceToInt(status[:]) == 0 {
			st = false
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

		var columnLen uint8
		buf := bytes.NewBuffer(columnLenBytes[:])
		err = binary.Read(buf, binary.LittleEndian, &columnLen)
		if err != nil {
			panic(err)
		}
		// 6. 读 ColumnData
		colData := make([]byte, columnLen)
		_, err = db.FileHandler.ReadAt(colData, offset)
		if err != nil {
			panic(err)
		}

		found := false
		column := metaDB.FindColumnByID(utils.SliceToUUID(tableID), utils.SliceToUUID(columnID))
		if column == nil {
			fmt.Println("tableID = ", utils.SliceToString(tableID))
			fmt.Println("columnID = ", utils.SliceToString(columnID))
			panic("Column is nil")
		}
		for _, td := range db.TbData {
			if td.TableID == utils.SliceToUUID(tableID) && td.Data["ID"].Value == utils.SliceToString(dataID) {
				found = true
				fmt.Println("tableID = ", utils.SliceToString(tableID[:]), "ColumnName = ", utils.SliceToString(column.Name[:]), "Offset = ", offset)
				td.Data[utils.SliceToString(column.Name[:])] = &ColumnData{
					PKValue: utils.SliceToString(dataID),
					Value:   utils.SliceToString(colData),
					Offset:  offset,
					Length:  int(columnLen),
					Status:  st,
				}
				break
			}
		}
		if !found {
			data := make(map[string]*ColumnData)
			//data["ID"] = utils.SliceToString(dataID)
			fmt.Println("tableID = ", utils.SliceToString(tableID[:]), "ColumnName = ", utils.SliceToString(column.Name[:]), "Offset = ", offset)
			data[utils.SliceToString(column.Name[:])] = &ColumnData{
				PKValue: utils.SliceToString(dataID),
				Value:   utils.SliceToString(colData),
				Offset:  offset,
				Length:  int(columnLen),
				Status:  st,
			}
			tbData := &TableData{
				TableID: utils.SliceToUUID(tableID),
				Data:    data,
			}
			db.TbData = append(db.TbData, tbData)
		}

		size -= int64(columnLen)
		offset += int64(columnLen)
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

func (db *DataDB) AddData(tableID [32]byte, dataID [32]byte, columns []*ColumnMetaData, model interface{}) {
	// 数据文件结构：TableID + DataID + Status + ColumnID + ColumnLength + ColumnData
	// 各个字段字节数分别为：32 + 32 + 1 + 32 + 1 + [ColumnLen]
	var buf bytes.Buffer
	data := make(map[string]*ColumnData)
	offset, _ := db.FileHandler.Seek(0, io.SeekEnd)
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
		st := true
		if utils.SliceToInt(status[:]) == 0 {
			st = false
		}
		err = binary.Write(&buf, binary.LittleEndian, column.ID)
		if err != nil {
			panic(err)
		}
		var columnLength = [1]byte{column.Length}
		err = binary.Write(&buf, binary.LittleEndian, columnLength)
		if err != nil {
			panic(err)
		}
		offset += 32 + 32 + 1 + 32 + 1
		//offset += int64(len(buf.Bytes()))
		fmt.Println("tableID = ", utils.SliceToString(tableID[:]), "columnName = ", utils.SliceToString(column.Name[:]), "offset = ", offset)
		for _, field := range GetModelFields(model) {
			if field.GetName() != utils.SliceToString(column.Name[:]) {
				continue
			}
			err = binary.Write(&buf, binary.LittleEndian, field.Bytes())
			if err != nil {
				panic(err)
			}
			data[utils.SliceToString(column.Name[:])] = &ColumnData{
				PKValue: utils.UUIDToString(dataID),
				Value:   field.GetValue(),
				Offset:  offset,
				Length:  int(column.Length),
				Status:  st,
			}
			break
		}
		offset += int64(column.Length)
	}
	err := binary.Write(db.FileHandler, binary.LittleEndian, buf.Bytes())
	if err != nil {
		panic(err)
	}
	d := &TableData{
		TableID: tableID,
		Data:    data,
	}
	db.TbData = append(db.TbData, d)
}

func (db *DataDB) SelectOneData(tableID [32]byte, conditions []Condition) map[string]string {
	fmt.Println("querySlice = ", conditions)
	for _, d := range db.TbData {
		if d.TableID != tableID {
			continue
		}
		flag := true
		for _, cond := range conditions {
			status := CompareByOperator(d.Data[cond.Field].Value, cond.Value, cond.Operator)
			if !status {
				flag = false
				break
			}
			fmt.Println("Yes, Equal !")
		}
		if flag && d.Data["ID"].Status == true {
			return ColumnDataToMap(d.Data)
		}
	}
	return nil
}

func (db *DataDB) SelectData(tableID [32]byte, conditions []Condition) []map[string]string {
	var result []map[string]string
	fmt.Println("querySlice = ", conditions)
	for _, d := range db.TbData {
		if d.TableID != tableID {
			continue
		}
		fmt.Println("dataMap = ", d.Data)
		flag := true
		for _, cond := range conditions {
			status := CompareByOperator(d.Data[cond.Field].Value, cond.Value, cond.Operator)
			if !status {
				flag = false
				break
			}
			fmt.Println("Yes, Equal !")
		}
		if flag && d.Data["ID"].Status == true {
			result = append(result, ColumnDataToMap(d.Data))
		}
	}
	return result
}

func (db *DataDB) UpdateData(tableID [32]byte, conditions []Condition, data map[string]string) {
	var result []map[string]string
	for _, d := range db.TbData {
		if d.TableID != tableID {
			continue
		}
		fmt.Println("dataMap = ", d.Data)
		flag := true
		for _, cond := range conditions {
			status := CompareByOperator(d.Data[cond.Field].Value, cond.Value, cond.Operator)
			if !status {
				flag = false
				break
			}
			fmt.Println("Yes, Equal !")
		}
		if flag && d.Data["ID"].Status == true {
			for key, value := range data {
				if _, found := d.Data[key]; found {
					d.Data[key].Value = value
					_, err := db.FileHandler.Seek(d.Data[key].Offset, io.SeekStart)
					if err != nil {
						panic(err)
					}
					fmt.Println("Update key = ", key, "value = ", value, "offset = ", d.Data[key].Offset)
					y := make([]byte, d.Data[key].Length)
					copy(y[:], value)
					err = binary.Write(db.FileHandler, binary.LittleEndian, y)
					if err != nil {
						panic(err)
					}
				}
			}
			result = append(result, ColumnDataToMap(d.Data))
		}
	}
}

func (db *DataDB) DeleteData(tableID [32]byte, conditions []Condition) {
	for _, d := range db.TbData {
		if d.TableID != tableID {
			continue
		}
		fmt.Println("dataMap = ", d.Data)
		flag := true
		for _, cond := range conditions {
			status := CompareByOperator(d.Data[cond.Field].Value, cond.Value, cond.Operator)
			if !status {
				flag = false
				break
			}
			fmt.Println("Yes, Equal !")
		}
		if flag && d.Data["ID"].Status == true {
			for _, d := range d.Data {
				// 获取 status 的 offset
				offset := d.Offset - 41
				_, err := db.FileHandler.Seek(offset, io.SeekStart)
				if err != nil {
					panic(err)
				}
				var status = [1]byte{0}
				err = binary.Write(db.FileHandler, binary.LittleEndian, status)
				if err != nil {
					panic(err)
				}
				d.Status = false
			}
		}
	}
}

func (db *MetaDB) CreateDB(name string) *DBMetaData {
	r := db.FindDBByName(name)
	if r != nil {
		return r
	}
	uid := utils.NewUUID()
	var nm [50]byte
	copy(nm[:], name)
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
		if column.TbID == tableID && utils.SliceToString(column.Name[:]) == columnName {
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
