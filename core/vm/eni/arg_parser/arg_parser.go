// type format (type_info) grammar
// type_info describe the types with encoding
// type: bool | int | uint | address | bytes | enum | string | fix_array | dyn_array | struct
// fix_array: fix_array_start [0-9]+ type
// dyn_array: dyn_array_start type
// struct: struct_start type+ struct_end

// TODO: string, bytes
package arg_parser

import (
	"bytes"
	"fmt"
	"math/big"
)

// token constant
const (
	BOOL = iota
	ADDRESS
	BYTES
	ENUM
	STRING
	FIX_ARRAY_START
	DYN_ARRAY_START
	STRUCT_START
	STRUCT_END
	INT
	INT8
	INT16
	INT24
	INT32
	INT40
	INT48
	INT56
	INT64
	INT72
	INT80
	INT88
	INT96
	INT104
	INT112
	INT120
	INT128
	INT136
	INT144
	INT152
	INT160
	INT168
	INT176
	INT184
	INT192
	INT200
	INT208
	INT216
	INT224
	INT232
	INT240
	INT248
	INT256
	UINT
	UINT8
	UINT16
	UINT24
	UINT32
	UINT40
	UINT48
	UINT56
	UINT64
	UINT72
	UINT80
	UINT88
	UINT96
	UINT104
	UINT112
	UINT120
	UINT128
	UINT136
	UINT144
	UINT152
	UINT160
	UINT168
	UINT176
	UINT184
	UINT192
	UINT200
	UINT208
	UINT216
	UINT224
	UINT232
	UINT240
	UINT248
	UINT256
	BYTE1
	BYTE2
	BYTE3
	BYTE4
	BYTE5
	BYTE6
	BYTE7
	BYTE8
	BYTE9
	BYTE10
	BYTE11
	BYTE12
	BYTE13
	BYTE14
	BYTE15
	BYTE16
	BYTE17
	BYTE18
	BYTE19
	BYTE20
	BYTE21
	BYTE22
	BYTE23
	BYTE24
	BYTE25
	BYTE26
	BYTE27
	BYTE28
	BYTE29
	BYTE30
	BYTE31
	BYTE32
)

// need type parsing
var complexType = map[byte]bool{
	FIX_ARRAY_START: true,
	DYN_ARRAY_START: true,
	STRUCT_START:    true,
	STRING:          true,
}

// in bytes
// only for value type
var dataLen = map[byte]int{
	BOOL:    1,
	INT:     32,
	UINT:    32,
	INT8:    1,
	INT16:   2,
	INT24:   3,
	INT32:   4,
	INT40:   5,
	INT48:   6,
	INT56:   7,
	INT64:   8,
	INT72:   9,
	INT80:   10,
	INT88:   11,
	INT96:   12,
	INT104:  13,
	INT112:  14,
	INT120:  15,
	INT128:  16,
	INT136:  17,
	INT144:  18,
	INT152:  19,
	INT160:  20,
	INT168:  21,
	INT176:  22,
	INT184:  23,
	INT192:  24,
	INT200:  25,
	INT208:  26,
	INT216:  27,
	INT224:  28,
	INT232:  29,
	INT240:  30,
	INT248:  31,
	INT256:  32,
	UINT8:   1,
	UINT16:  2,
	UINT24:  3,
	UINT32:  4,
	UINT40:  5,
	UINT48:  6,
	UINT56:  7,
	UINT64:  8,
	UINT72:  9,
	UINT80:  10,
	UINT88:  11,
	UINT96:  12,
	UINT104: 13,
	UINT112: 14,
	UINT120: 15,
	UINT128: 16,
	UINT136: 17,
	UINT144: 18,
	UINT152: 19,
	UINT160: 20,
	UINT168: 21,
	UINT176: 22,
	UINT184: 23,
	UINT192: 24,
	UINT200: 25,
	UINT208: 26,
	UINT216: 27,
	UINT224: 28,
	UINT232: 29,
	UINT240: 30,
	UINT248: 31,
	UINT256: 32,
	BYTE1:   1,
	BYTE2:   2,
	BYTE3:   3,
	BYTE4:   4,
	BYTE5:   5,
	BYTE6:   6,
	BYTE7:   7,
	BYTE8:   8,
	BYTE9:   9,
	BYTE10:  10,
	BYTE11:  11,
	BYTE12:  12,
	BYTE13:  13,
	BYTE14:  14,
	BYTE15:  15,
	BYTE16:  16,
	BYTE17:  17,
	BYTE18:  18,
	BYTE19:  19,
	BYTE20:  20,
	BYTE21:  21,
	BYTE22:  22,
	BYTE23:  23,
	BYTE24:  24,
	BYTE25:  25,
	BYTE26:  26,
	BYTE27:  27,
	BYTE28:  28,
	BYTE29:  29,
	BYTE30:  30,
	BYTE31:  31,
	BYTE32:  32}

func Parse(type_info []byte, data []byte) string {
	var json bytes.Buffer
	parse_entry_point(type_info, data, &json)
	return json.String()
}

func parse_entry_point(type_info []byte, data []byte, json *bytes.Buffer) {
	json.WriteString("[")
	for i := 0; 0 < len(type_info); i++ {
		if 0 < i {
			json.WriteString(",")
		}
		type_info, data = parse_type(type_info, data, json)
	}
	json.WriteString("]")
}

// assuming that data are packed
func parse_type(type_info []byte, data []byte, json *bytes.Buffer) ([]byte, []byte) {
	t := type_info[0]
	if complexType[t] {
		if t == FIX_ARRAY_START {
			type_info, data = parse_fix_array(type_info, data, json)
		} else if t == DYN_ARRAY_START {
			type_info, data = parse_dyn_array(type_info, data, json)
		} else if t == STRUCT_START {
			type_info, data = parse_struct(type_info, data, json)
		} else if t == STRING {
			type_info, data = parse_string(type_info, data, json)
		} else { // error

		}
	} else { // value type
		type_info, data = parse_value(type_info, data, json)
	}
	return type_info, data
}

func parse_string(type_info []byte, data []byte, json *bytes.Buffer) ([]byte, []byte) {
	type_info = type_info[1:] // string
	leng := new(big.Int).SetBytes(data[:32]).Int64()
	data = data[32:]

	var buffer bytes.Buffer
	for i := int64(0); i < leng; i++ {
		if data[i] == '\\' || data[i] == '"' {
			buffer.WriteByte('\\')
		}
		buffer.WriteByte(data[i])
	}
	json.WriteString("\"")
	json.WriteString(buffer.String())
	json.WriteString("\"")
	data = data[leng:]
	if leng%32 > 0 {
		data = data[32-leng%32:]
	}
	return type_info, data
}

func parse_fix_array(type_info []byte, data []byte, json *bytes.Buffer) ([]byte, []byte){
    type_info = type_info[1:] // fix_array_start
    json.WriteString("[")
    leng := new(big.Int).SetBytes(type_info[:32]).Int64()
    type_info = type_info[32:]

    for i:=int64(0); i<leng; i++{
        if i==leng-1 {
            type_info, data = parse_type(type_info, data, json)
        }else{
            json.WriteString(", ")
            _, data = parse_type(type_info, data, json)
        }
    }

	json.WriteString("]")
	return type_info, data
}

// dynamic array
func parse_dyn_array(type_info []byte, data []byte, json *bytes.Buffer) ([]byte, []byte) {
	// TODO
	return type_info, data
}

func parse_struct(type_info []byte, data []byte, json *bytes.Buffer) ([]byte, []byte) {
	type_info = type_info[1:] // struct_start
	json.WriteString("[")
	for i := 0; ; i++ {
		t := type_info[0]
		if 0 < i {
			json.WriteString(", ")
		}
		if t == STRUCT_END {
			break
		}
		type_info, data = parse_type(type_info, data, json)
	}
	type_info = type_info[1:] // struct_end
	json.WriteString("]")
	return type_info, data
}

// bool, int
func parse_value(type_info []byte, data []byte, json *bytes.Buffer) ([]byte, []byte) {
	t := type_info[0]
	if t == BOOL {
		json.WriteString(fmt.Sprint(0 != data[0]))
	} else if INT <= t && t <= INT256 { // signed integer
		n := new(big.Int)
		var b [32]byte
		copy(b[:], data[:dataLen[t]])
		if b[0] >= 128 { // negative value, two's complement
			n.SetBytes(b[:])
			n = n.Sub(n, big.NewInt(int64(1)))
			copy(b[:], n.Bytes())
			for i := 0; i < 32; i++ {
				b[i] ^= 255
			}
			n.SetBytes(b[:])
			n = n.Mul(n, big.NewInt(int64(-1)))
			json.WriteString(n.String())
		} else { // positive value
			n.SetBytes(b[:])
			json.WriteString(n.String())
		}

	} else if (UINT <= t && t <= UINT256) || (BYTE1 <= t && t <= BYTE32) { // unsigned integer
		n := new(big.Int)
		n.SetBytes(data[:dataLen[t]]) // big endian
		json.WriteString(n.String())
	}
	type_info = type_info[1:]
	data = data[dataLen[t]:]
	return type_info, data
}
