package redis

import "github.com/codecrafters-io/redis-starter-go/resp"

type StorageType string

const (
	KeyStorage    StorageType = "key"
	ListStorage   StorageType = "list"
	StreamStorage StorageType = "stream"
)

type StorageField struct {
	Type   StorageType
	Key    *resp.RESPValue
	List   []*resp.RESPValue
	Stream *StreamField
}

type StreamField struct {
	Entries   map[int64]map[int64][]*resp.RESPValue
	LastEntry *StreamEntryID
}

type StreamEntryID struct {
	Time     int64
	Sequence int64
}

func (s *StorageField) IsKey() bool {
	return s.Type == KeyStorage
}

func (s *StorageField) IsList() bool {
	return s.Type == ListStorage
}

func (s *StorageField) IsStream() bool {
	return s.Type == StreamStorage
}
