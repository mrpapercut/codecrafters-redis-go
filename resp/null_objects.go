package resp

func NullBulkstring() *RESPValue {
	return &RESPValue{
		Type:   BulkString,
		IsNull: true,
	}
}

func NullArray() *RESPValue {
	return &RESPValue{
		Type:   Array,
		IsNull: true,
	}
}
