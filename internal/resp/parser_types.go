package resp

type Request struct {
	Command string
	Args    []string
}

type RespType int

const (
	SimpleString RespType = iota
	Error
	Integer
	BulkString
	Array
)

type RespValue struct {
	Type RespType

	String  string
	Integer int64
	Array   []RespValue

	// expiry for TTL
	ExpiresAt int64

	// only true if BulkString or Array is null
	Null bool
}

func NewError(s string) RespValue {
	return RespValue{
		Type:   Error,
		String: s,
	}
}

func NewSimpleString(s string) RespValue {
	return RespValue{
		Type:   SimpleString,
		String: s,
	}
}

func NewBulkString(s string) RespValue {
	return RespValue{
		Type:   BulkString,
		String: s,
	}
}

func NewInteger(i int64) RespValue {
	return RespValue{
		Type:    Integer,
		Integer: i,
	}
}

func NewNullBulkString() RespValue {
	return RespValue{
		Type: BulkString,
		Null: true,
	}
}

func NewNullArr() RespValue {
	return RespValue{
		Type: Array,
		Null: true,
	}
}
