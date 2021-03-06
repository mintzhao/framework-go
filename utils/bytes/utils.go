package bytes

import (
	"encoding/binary"
	"math"
)

/**
 * @Author: imuge
 * @Date: 2020/4/28 5:42 下午
 */

var (
	DEFAULT_CHARSET = "UTF-8"
	MAX_BUFFER_SIZE = 1024 * 1024 * 1024
	BUFFER_SIZE     = 64
	TRUE_BYTE       = byte(0x01)
	FALSE_BYTE      = byte(0x00)
)

/**
 * 将 short 值转为2字节的二进制数组；
 *
 * @param value 要转换的int整数；
 * @param bytes 要保存转换结果的二进制数组；转换结果将从高位至低位的顺序写入数组从 0 开始的2个元素；
 */
func int16ToBytes(value int16, bytes []byte, offset int) int {
	bytes[offset] = (byte)((value >> 8) & 0x00FF)
	bytes[offset+1] = (byte)(value & 0x00FF)
	return 2
}

/**
 * 将 int 值转为4字节的二进制数组
 *
 * 以“高位在前”的方式转换，即：数值的高位保存在数组地址的低位；
 *
 * @param value  要转换的int整数；
 * @param bytes  要保存转换结果的二进制数组；转换结果将从高位至低位的顺序写入数组从 offset 指定位置开始的4个元素；
 * @param offset 写入转换结果的起始位置；
 * @return 返回写入的长度；
 */
func int32ToBytes(value int32, bytes []byte, offset int) int {
	bytes[offset] = (byte)((value >> 24) & 0x00FF)
	bytes[offset+1] = (byte)((value >> 16) & 0x00FF)
	bytes[offset+2] = (byte)((value >> 8) & 0x00FF)
	bytes[offset+3] = (byte)(value & 0x00FF)
	return 4
}

/**
 * 将 long 值转为8字节的二进制数组；
 *
 * @param value  要转换的long整数；
 * @param bytes  要保存转换结果的二进制数组；转换结果将从高位至低位的顺序写入数组从 offset 指定位置开始的8个元素；
 * @param offset 写入转换结果的起始位置；
 * @return 返回写入的长度；
 */
func int64ToBytes(value int64, bytes []byte, offset int) int {
	bytes[offset] = (byte)((value >> 56) & 0x00FF)
	bytes[offset+1] = (byte)((value >> 48) & 0x00FF)
	bytes[offset+2] = (byte)((value >> 40) & 0x00FF)
	bytes[offset+3] = (byte)((value >> 32) & 0x00FF)
	bytes[offset+4] = (byte)((value >> 24) & 0x00FF)
	bytes[offset+5] = (byte)((value >> 16) & 0x00FF)
	bytes[offset+6] = (byte)((value >> 8) & 0x00FF)
	bytes[offset+7] = (byte)(value & 0x00FF)
	return 8
}

/**
 * 将 int 值转为4字节的二进制数组；
 *
 * @param value value
 * @return 转换后的二进制数组，高位在前，低位在后；
 */
func Int32ToBytes(value int32) []byte {
	bytes := make([]byte, 4)
	int32ToBytes(value, bytes, 0)
	return bytes
}

func IntToBytes(value int) []byte {
	if value <= math.MaxInt32 {
		return Int32ToBytes(int32(value))
	} else {
		return Int64ToBytes(int64(value))
	}
}

func ToInt(bys []byte) int {
	if len(bys) <= 4 {
		return int(ToInt32(bys))
	} else {
		return int(ToInt64(bys))
	}
}

func BoolToBytes(value bool) byte {
	if value {
		return TRUE_BYTE
	} else {
		return FALSE_BYTE
	}
}

/**
 * 将 long 值转为8字节的二进制数组；
 *
 * @param value value
 * @return 转换后的二进制数组，高位在前，低位在后；
 */
func Int64ToBytes(value int64) []byte {
	bytes := make([]byte, 8)
	int64ToBytes(value, bytes, 0)
	return bytes
}

func Int16ToBytes(value int16) []byte {
	bytes := make([]byte, 2)
	int16ToBytes(value, bytes, 0)
	return bytes
}

func Int8ToBytes(value int8) byte {
	return (byte)(value)
}

// TODO UTF-8 ?
func StringToBytes(str string) []byte {
	return []byte(str)
}

func ToString(bytes []byte) string {
	return string(bytes)
}

func ToBoolean(b byte) bool {
	if b == TRUE_BYTE {
		return true
	}

	return false
}

func ToInt8(b byte) int8 {
	return int8(b)
}

func ToInt16(b []byte) int16 {
	if len(b) < 2 {
		for i := 0; i < 2-len(b); i++ {
			b = append([]byte{0x00}, b...)
		}
	}
	return int16(binary.BigEndian.Uint16(b))
}

func ToInt32(b []byte) int32 {
	if len(b) < 4 {
		for i := 0; i < 4-len(b); i++ {
			b = append([]byte{0x00}, b...)
		}
	}
	return int32(binary.BigEndian.Uint32(b))
}

func ToInt64(b []byte) int64 {
	if len(b) < 8 {
		for i := 0; i < 8-len(b); i++ {
			b = append([]byte{0x00}, b...)
		}
	}
	return int64(binary.BigEndian.Uint64(b))
}

func Concat(slices ...[]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func StartsWith(srcBytes []byte, prefixBytes []byte) bool {
	for i := 0; i < len(prefixBytes); i++ {
		if prefixBytes[i] != srcBytes[i] {
			return false
		}
	}
	return true
}

/**
 * 比较指定的两个字节数组是否一致；
 * 此方法不处理两者其中之一为 nil 的情形
 */
func Equals(bytes1 []byte, bytes2 []byte) bool {
	if string(bytes1) == string(bytes2) {
		return true
	}
	if bytes1 == nil || bytes2 == nil {
		return false
	}
	if len(bytes1) != len(bytes2) {
		return false
	}
	for i := 0; i < len(bytes1); i++ {
		if bytes1[i] != bytes2[i] {
			return false
		}
	}
	return true
}
