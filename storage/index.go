package storage

// Index 索引
type Index struct {
	FileID  [32]byte
	BlockID int8
	Offset  int32
	Size    int32
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
		BlockID: int8(blockID),
		Offset:  int32(offset),
		Size:    int32(size),
	}
	return &index
}
