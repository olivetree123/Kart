package storage

// Index 索引
type Index struct {
	FileID  [32]byte // 32 byte
	Offset  int64    // 8 byte
	Size    int64    // 8 byte
	BlockID int64    // 8 byte
}

// NewIndex 创建索引
func NewIndex(fileID string, blockID int, offset int, size int) *Index {
	if len(fileID) > 32 {
		panic("fileID too long.")
	}
	var fID [32]byte
	copy(fID[:], fileID)
	index := Index{
		FileID:  fID,
		BlockID: int64(blockID),
		Offset:  int64(offset),
		Size:    int64(size),
	}
	return &index
}
