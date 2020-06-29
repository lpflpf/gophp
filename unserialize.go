package gophp

import (
	"bytes"
	"errors"
	"fmt"
)

var ErrNotExpectedType = errors.New("UnMarshal: Not expect type")
var ErrNotExpectedToken = errors.New("UnMarshal: Not expect token")
var ErrReadLengthFailed = errors.New("UnMarshal: Error while reading length of value")

func UnMarshal(data []byte) (interface{}, error) {
	reader := bytes.NewReader(data)
	return unMarshalByReader(reader)
}

func unMarshalByReader(reader *bytes.Reader) (interface{}, error) {
	for {
		if token, _, err := reader.ReadRune(); err == nil {
			switch token {
			default:
				return nil, ErrNotExpectedToken
			case 'N':
				return unMarshalNil(reader)
			case 'b':
				return unMarshalBool(reader)
			case 'i':
				return unMarshalInt64(reader)
			case 'd':
				return unMarshalFloat64(reader)
			case 's':
				return unMarshalString(reader, true)
			case 'a':
				return unMarshalArray(reader)
			case 'O':
				return unMarshalObject(reader)
			case 'R', 'r': // R 带指针的引用，r不带指针的引用
				// 忽略
			case 'C':
				// 自定义序列化
			}
		}
		return nil, nil
	}

}

func unMarshalNil(reader *bytes.Reader) (interface{}, error) {
	return nil, expect(reader, ';')
}

func unMarshalBool(reader *bytes.Reader) (interface{}, error) {
	var (
		raw rune
		err error
	)
	err = expect(reader, ':')
	if err != nil {
		return nil, err
	}

	if raw, _, err = reader.ReadRune(); err != nil {
		return nil, ErrNotExpectedToken
	}

	err = expect(reader, ';')
	if err != nil {
		return nil, err
	}
	return raw == '1', nil
}

func unMarshalInt64(reader *bytes.Reader) (interface{}, error) {
	var ret int64
	if n, err := fmt.Fscanf(reader, ":%d;", &ret); err != nil || n == 0 {
		return nil, err
	}
	return ret, nil
}

func unMarshalFloat64(reader *bytes.Reader) (interface{}, error) {
	var ret float64
	if n, err := fmt.Fscanf(reader, ":%lf;", &ret); err != nil || n == 0 {
		return nil, err
	}
	return ret, nil
}

func unMarshalString(reader *bytes.Reader, isFinal bool) (interface{}, error) {
	var (
		err     error
		val     interface{}
		strLen  int
		readLen int
	)

	strLen, err = readLength(reader)

	err = expect(reader, '"')
	if err != nil {
		return nil, err
	}

	if strLen > 0 {
		buf := make([]byte, strLen, strLen)
		if readLen, err = reader.Read(buf); err != nil {
			return nil, ErrNotExpectedToken
		} else {
			if readLen != strLen {
				return nil, ErrNotExpectedToken
			} else {
				val = string(buf)
			}
		}
	}

	err = expect(reader, '"')
	if err != nil {
		return nil, err
	}
	if isFinal {
		err = expect(reader, ';')
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}

func unMarshalObject(reader *bytes.Reader) (interface{}, error) {
	_, _ = unMarshalString(reader, false)
	return unMarshalArray(reader)
}

func unMarshalArray(reader *bytes.Reader) (interface{}, error) {
	var arrLen int
	var err error
	val := make(map[string]interface{})
	var val2 []interface{}

	arrLen, err = readLength(reader)

	if err != nil {
		return nil, err
	}
	err = expect(reader, '{')
	if err != nil {
		return nil, err
	}

	indexLen := 0
	for i := 0; i < arrLen; i++ {
		k, err := unMarshalByReader(reader)
		if err != nil {
			return nil, err
		}
		v, err := unMarshalByReader(reader)
		if err != nil {
			return nil, err
		}

		switch k.(type) {
		case string, float64:
			stringKey, _ := k.(string)
			val[stringKey] = v
			val2 = append(val2, v)
		case int64:
			stringKey, _ := NumericalToString(k)
			val[stringKey] = v
			val2 = append(val2, v)

			if int64(i) == k.(int64) {
				indexLen++
			}
		default:
			return nil, ErrNotExpectedType
		}
	}

	err = expect(reader, '}')
	if err != nil {
		return nil, err
	}

	if indexLen == arrLen {
		return val2, nil
	}

	return val, nil
}

func expect(reader *bytes.Reader, expected rune) error {
	if token, _, err := reader.ReadRune(); err != nil || token != expected {
		return ErrNotExpectedToken
	}
	return nil
}

func readLength(reader *bytes.Reader) (ret int, _ error) {
	if n, err := fmt.Fscanf(reader, ":%d:", &ret); n == 0 || err != nil {
		return 0, ErrReadLengthFailed
	}
	return
}
