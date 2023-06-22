package nbt

import (
	"encoding/binary"
	"errors"
	"math"
)

type decoder struct {
	data   []byte
	cursor int
}

func Decode(data []byte) (map[string]*NbtTag, error) {
	d := &decoder{
		data: data,
	}

	return d.readCompound(true)
}

func (d *decoder) readBytes(n int) []byte {
	if len(d.data)-d.cursor < n {
		return nil
	}

	d.cursor += n
	return d.data[d.cursor-n : d.cursor]
}

func (d *decoder) readCompound(singleCompound bool) (map[string]*NbtTag, error) {
	tags := make(map[string]*NbtTag)
	for {
		name, tag, err := d.readTag()
		if err != nil {
			return nil, err
		}

		if tag.Type == TAG_End {
			break
		}

		tags[name] = tag

		if singleCompound == true {
			break
		}
	}

	return tags, nil
}

func (d *decoder) readTag() (string, *NbtTag, error) {
	tag := &NbtTag{}

	typeIdBytes := d.readBytes(1)
	if typeIdBytes == nil {
		tag.Type = TAG_End
		return "", tag, nil
	}

	typeId := typeIdBytes[0]
	tag.Type = typeId

	if typeId == TAG_End {
		return "", tag, nil
	}

	name, err := d.readString()
	if err != nil {
		return "", nil, err
	}

	val, err := d.readTagValue(typeId)
	if err != nil {
		return "", nil, err
	}
	if val == nil {
		tag.Type = TAG_End
		return "", tag, nil
	}
	tag.Value = val

	return name, tag, nil
}

func (d *decoder) readTagValue(typeId uint8) (any, error) {
	switch typeId {
	case TAG_Byte:
		return d.readByte()
	case TAG_Short:
		return d.readShort()
	case TAG_Int:
		return d.readInt()
	case TAG_Long:
		return d.readLong()
	case TAG_Float:
		return d.readFloat()
	case TAG_Double:
		return d.readDouble()
	case TAG_Byte_Array:
		return d.readByteArray()
	case TAG_String:
		return d.readString()
	case TAG_List:
		return d.readList()
	case TAG_Compound:
		return d.readCompound(false)
	case TAG_End:
		return nil, nil
	case TAG_Int_Array:
		return d.readIntArray()
	case TAG_Long_Array:
		return d.readLongArray()
	}

	return nil, nil
}

func (d *decoder) readString() (string, error) {
	length, err := d.readUnsignedShort()
	if err != nil {
		return "", err
	}

	strBytes := d.readBytes(int(length))
	if strBytes == nil {
		return "", errors.New("Unexpected EOF")
	}

	return string(strBytes), nil
}

func (d *decoder) readByte() (int8, error) {
	bytes := d.readBytes(1)
	if bytes == nil {
		return 0, errors.New("Unexpected EOF")
	}
	return int8(bytes[0]), nil
}

func (d *decoder) readShort() (int16, error) {
	bytes := d.readBytes(2)
	if bytes == nil {
		return 0, errors.New("Unexpected EOF")
	}
	return int16(binary.BigEndian.Uint16(bytes)), nil
}

func (d *decoder) readInt() (int32, error) {
	bytes := d.readBytes(4)
	if bytes == nil {
		return 0, errors.New("Unexpected EOF")
	}
	return int32(binary.BigEndian.Uint32(bytes)), nil
}

func (d *decoder) readLong() (int64, error) {
	bytes := d.readBytes(8)
	if bytes == nil {
		return 0, errors.New("Unexpected EOF")
	}
	return int64(binary.BigEndian.Uint64(bytes)), nil
}

func (d *decoder) readFloat() (float32, error) {
	bytes := d.readBytes(4)
	if bytes == nil {
		return 0, errors.New("Unexpected EOF")
	}
	return math.Float32frombits(binary.BigEndian.Uint32(bytes)), nil
}

func (d *decoder) readDouble() (float64, error) {
	bytes := d.readBytes(8)
	if bytes == nil {
		return 0, errors.New("Unexpected EOF")
	}
	return math.Float64frombits(binary.BigEndian.Uint64(bytes)), nil
}

func (d *decoder) readByteArray() ([]int8, error) {
	length, err := d.readInt()
	if err != nil {
		return nil, err
	}

	bytes := make([]int8, length)
	for i := int32(0); i < length; i++ {
		bytes[i], err = d.readByte()
		if err != nil {
			return nil, err
		}
	}
	return bytes, nil
}

func (d *decoder) readList() ([]any, error) {
	typeIdBytes := d.readBytes(1)
	if typeIdBytes == nil {
		return nil, errors.New("Unexpected EOF")
	}
	typeId := typeIdBytes[0]

	length, err := d.readInt()
	if err != nil {
		return nil, err
	}
	if length <= 0 {
		return nil, nil
	}

	list := make([]any, length)
	for i := int32(0); i < length; i++ {
		list[i], err = d.readTagValue(typeId)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (d *decoder) readIntArray() ([]int32, error) {
	length, err := d.readInt()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, errors.New("Invalid array length")
	}

	array := make([]int32, length)
	for i := int32(0); i < length; i++ {
		array[i], err = d.readInt()
		if err != nil {
			return nil, err
		}
	}

	return array, nil
}

func (d *decoder) readLongArray() ([]int64, error) {
	length, err := d.readInt()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, errors.New("Invalid array length")
	}

	array := make([]int64, length)
	for i := int32(0); i < length; i++ {
		array[i], err = d.readLong()
		if err != nil {
			return nil, err
		}
	}

	return array, nil
}

func (d *decoder) readUnsignedShort() (uint16, error) {
	bytes := d.readBytes(2)
	if bytes == nil {
		return 0, errors.New("Unexpected EOF")
	}
	return binary.BigEndian.Uint16(bytes), nil
}
