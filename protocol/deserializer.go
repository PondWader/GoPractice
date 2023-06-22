package protocol

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

func (client *ProtocolClient) deserialize(data []byte, format interface{}) error {
	value := reflect.ValueOf(format).Elem()
	formatType := reflect.TypeOf(format).Elem()

	offset := 0
	deserializer := client.newDeserializer(func(l int) ([]byte, error) {
		if len(data) < offset+l {
			return nil, errors.New("EOF")
		}
		d := data[offset : offset+l]
		offset += l
		return d, nil
	})

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		valueType := formatType.Field(i).Tag.Get("type")

		switch valueType {
		case "VarInt":
			v, _, err := deserializer.readVarInt()
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v))

		case "String":
			v, err := deserializer.readString()
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v))

		case "Byte":
			v, err := deserializer.readBytes(1)
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(int8(v[0])))

		case "UnsignedShort":
			v, err := deserializer.readUnsignedShort()
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v))

		case "Long":
			v, err := deserializer.readLong()
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v))

		case "UnsignedLong":
			v, err := deserializer.readUnsignedLong()
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v))

		case "Position":
			v, err := deserializer.readUnsignedLong()
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}

			x := int32(v >> 38)
			y := int32((v >> 26) & 0xFFF)
			z := int32((v << 38) >> 38)
			field.Set(reflect.ValueOf(&Position{x, y, z}))

		case "ByteArray":
			lengthField := formatType.Field(i).Tag.Get("length")
			v, err := deserializer.readBytes(value.FieldByName(lengthField).Interface().(int))
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v))

		case "Double":
			v, err := deserializer.readDouble()
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v))

		case "Float":
			v, err := deserializer.readFloat()
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v))

		case "Boolean":
			v, err := deserializer.readBytes(1)
			if err != nil {
				client.Disconnect(err.Error())
				return err
			}
			field.Set(reflect.ValueOf(v[0] == 1))

		default:
			panic("A type of " + valueType + " is used but the deserializer does not know how to handle it!")
		}
	}

	return nil
}

type deserializer struct {
	readBytes func(int) ([]byte, error)
}

func (client *ProtocolClient) newDeserializer(readBytes func(int) ([]byte, error)) deserializer {
	return deserializer{readBytes}
}

func (deserializer *deserializer) readVarInt() (int, int, error) {
	value := 0
	position := 0
	SEGMENT_BITS := 0x7F
	CONTINUE_BIT := 0x80

	for {
		readBytes, err := deserializer.readBytes(1)
		if err != nil {
			return 0, 0, err
		}
		currentByte := readBytes[0]
		value |= (int(currentByte) & SEGMENT_BITS) << position

		if (int(currentByte) & CONTINUE_BIT) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, 0, errors.New("A VarInt with a size >= 32 was received")
		}
	}

	return value, (position + 7) / 7, nil
}

func (deserializer *deserializer) readString() (string, error) {
	length, _, err := deserializer.readVarInt()
	if err != nil {
		return "", err
	}
	bytes, err := deserializer.readBytes(length)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (deserializer *deserializer) readUnsignedShort() (uint16, error) {
	bytes, err := deserializer.readBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(bytes), nil
}

func (deserializer *deserializer) readLong() (int64, error) {
	bytes, err := deserializer.readBytes(8)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(bytes)), nil
}

func (deserializer *deserializer) readUnsignedLong() (uint64, error) {
	bytes, err := deserializer.readBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(bytes), nil
}

func (deserializer *deserializer) readDouble() (float64, error) {
	bytes, err := deserializer.readBytes(8)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(binary.BigEndian.Uint64(bytes)), nil
}

func (deserializer *deserializer) readFloat() (float32, error) {
	bytes, err := deserializer.readBytes(4)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(binary.BigEndian.Uint32(bytes)), nil
}
