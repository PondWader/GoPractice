package protocol

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"reflect"

	"github.com/PondWader/GoPractice/utils"
	"github.com/google/uuid"
)

func Serialize(format interface{}) []byte {
	data := make([]byte, 0)

	reflectValue := reflect.ValueOf(format).Elem()
	formatType := reflect.TypeOf(format).Elem()

	for i := 0; i < reflectValue.NumField(); i++ {
		value := reflectValue.Field(i).Interface()
		valueType := formatType.Field(i).Tag.Get("type")

		valueIf := formatType.Field(i).Tag.Get("if")
		if valueIf != "" {
			if reflectValue.FieldByName(valueIf).Bool() == false {
				continue
			}
		}

		switch valueType {
		case "VarInt":
			data = append(data, writeVarInt(value.(int))...)
		case "String":
			data = append(data, writeString(value.(string))...)
		case "JSON":
			jsonVal, _ := json.Marshal(value)
			data = append(data, writeString(string(jsonVal))...)
		case "ByteArray":
			data = append(data, value.([]byte)...)
		case "Int":
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, uint32(value.(int32)))
			data = append(data, bytes...)
		case "Long":
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, uint64(value.(int64)))
			data = append(data, bytes...)
		case "UnsignedByte":
			data = append(data, value.(uint8))
		case "Byte":
			data = append(data, uint8(value.(int8)))
		case "UnsignedShort":
			bytes := make([]byte, 2)
			binary.BigEndian.PutUint16(bytes, value.(uint16))
			data = append(data, bytes...)
		case "Boolean":
			v := value.(bool)
			if v == true {
				data = append(data, 1)
			} else {
				data = append(data, 0)
			}
		case "Float":
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, math.Float32bits(value.(float32)))
			data = append(data, bytes...)
		case "Double":
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, math.Float64bits(value.(float64)))
			data = append(data, bytes...)
		case "UUID":
			bytes, _ := value.(*uuid.UUID).MarshalBinary()
			data = append(data, bytes...)
		case "Array":
			t := unpackArray(value)
			for _, item := range t {
				data = append(data, Serialize(item)...)
			}
		default:
			utils.Error("Cannot serialize value of type", valueType)
		}
	}

	return data
}

// https://stackoverflow.com/a/73029665
func unpackArray(s any) []any {
	v := reflect.ValueOf(s)
	r := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		r[i] = v.Index(i).Interface()
	}
	return r
}

func writeVarInt(value int) []byte {
	SEGMENT_BITS := 0x7F
	CONTINUE_BIT := 0x80

	bytes := []byte{}

	for {
		if (value & -128) == 0 {
			bytes = append(bytes, byte(value))
			return bytes
		}

		bytes = append(bytes, byte((value&SEGMENT_BITS)|CONTINUE_BIT))

		// Change to uint so can move sign bit
		unsigned := uint(value)
		unsigned >>= 7
		value = int(unsigned)
	}
}

func writeString(str string) []byte {
	return append(writeVarInt(len(str)), str...)
}
