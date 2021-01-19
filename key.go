package hierarchical

import (
	"bytes"
	"reflect"
	"unsafe"
)

type AppendMarshaler interface {
	AppendMarshal([]byte) ([]byte, error)
}

type BinaryKey struct {
	Id       []byte
	Property []byte
	Key      []byte
}

func ReadBinaryKey(data []byte) BinaryKey {
	var result BinaryKey
	index := bytes.IndexByte(data, 0)
	if index == -1 {
		return BinaryKey{Id: data}
	}
	result.Id = data[:index]
	data = data[index+1:]
	index = bytes.IndexByte(data, 0)
	if index == -1 {
		result.Property = data
		return result
	}
	result.Property = data[:index]
	result.Key = data[index+1:]
	return result
}

type Key struct {
	Id       string
	Property string
	Key      []byte
}

func ReadKey(data []byte) Key {
	result := ReadBinaryKey(data)
	return Key{string(result.Id), string(result.Property), result.Key}
}

func AppendKey(data []byte, id string, property string, key AppendMarshaler) ([]byte, error) {
	sheader := (*reflect.StringHeader)(unsafe.Pointer(&id))
	var header reflect.SliceHeader
	header.Data = sheader.Data
	header.Len = sheader.Len
	header.Cap = sheader.Len
	data = append(data, *((*[]byte)(unsafe.Pointer(&header)))...)
	sheader = (*reflect.StringHeader)(unsafe.Pointer(&property))
	header.Data = sheader.Data
	header.Len = sheader.Len
	header.Cap = sheader.Len
	data = append(data, 0)
	data = append(data, *((*[]byte)(unsafe.Pointer(&header)))...)
	if key == nil {
		return data, nil
	}
	data = append(data, 0)
	return key.AppendMarshal(data)
}

func AppendId(data []byte, id string) []byte {
	sheader := (*reflect.StringHeader)(unsafe.Pointer(&id))
	var header reflect.SliceHeader
	header.Data = sheader.Data
	header.Len = sheader.Len
	header.Cap = sheader.Len
	return append(data, *((*[]byte)(unsafe.Pointer(&header)))...)
}
