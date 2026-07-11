package resp

import (
	"strconv"
	"strings"
)

func Encode(v RespValue) string {
	switch v.Type {
	case SimpleString:
		return "+" + v.String + "\r\n"
	case Error:
		return "-" + v.String + "\r\n"
	case Integer:
		return ":" + strconv.FormatInt(v.Integer, 10) + "\r\n"
	case BulkString:
		if v.Null {
			return "$-1\r\n"
		}
		return "$" + strconv.Itoa(len(v.String)) + "\r\n" + v.String + "\r\n"
	case Array:
		if v.Null {
			return "*-1\r\n"
		}
		var b strings.Builder
		b.WriteString("*" + strconv.Itoa(len(v.Array)) + "\r\n")
		for _, elem := range v.Array {
			b.WriteString(Encode(elem))
		}
		return b.String()
	default:
		return ""
	}
}
