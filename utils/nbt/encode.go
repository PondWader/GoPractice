package nbt

import "encoding/binary"

func Encode(nbt map[string]*NbtTag) []byte {
	data := []byte{}
	encodeCompound(nbt, data)
	return data
}

func encodeCompound(nbt map[string]*NbtTag, data []byte) {
	for name, tag := range nbt {
		data = append(data, tag.Type)
		if tag.Type == TAG_End {
			return
		}

		writeString(data, name)

		switch tag.Type {
		case TAG_Compound:
			encodeCompound(tag.Value.(map[string]*NbtTag), data)
		case TAG_Byte:
			data = append(data, byte(tag.Value.(int8)))
		case TAG_String:
			writeString(data, tag.Value.(string))
		case TAG_Short:
			writeShort(data, tag.Value.(int16))
		}
	}
}

func writeString(data []byte, str string) {
	writeUnsignedShort(data, uint16(len(str)))
	data = append(data, []byte(str)...)
}

func writeShort(data []byte, n int16) {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(n))
	data = append(data, bytes...)
}

func writeUnsignedShort(data []byte, n uint16) {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, n)
	data = append(data, bytes...)
}
