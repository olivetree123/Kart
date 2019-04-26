package storage

// 缺少：文件名称、上传时间、更新时间、访问次数

// Index 索引
type Index struct {
	ID       [32]byte // 32 byte
	BucketID [32]byte
	Offset   int64 // 8 byte
	Size     int64 // 8 byte
	BlockID  int64 // 8 byte
}

// NewIndex 创建索引
func NewIndex(fileID string, bucketID string, blockID int, offset int, size int) *Index {
	if len(fileID) > 32 {
		panic("fileID too long.")
	}
	var fID, bID [32]byte
	copy(fID[:], fileID)
	copy(bID[:], bucketID)
	index := &Index{
		ID:       fID,
		BucketID: bID,
		BlockID:  int64(blockID),
		Offset:   int64(offset),
		Size:     int64(size),
	}
	return index
}
