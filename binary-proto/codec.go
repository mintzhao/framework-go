package binary_proto

import (
	"framework-go/utils/bytes"
	"reflect"
	"strconv"
)

/**
 * @Author: imuge
 * @Date: 2020/5/21 下午
 */

var Cdc = Codec{
	make(map[int32]interface{}),
	make(map[int32]interface{}),
	make(map[int32]interface{}),
}

type Codec struct {
	// 契约类型 // TODO 未支持多版本
	contractMap map[int32]interface{}
	// 枚举类型 // TODO 未支持多版本
	enumMap map[int32]interface{}
	// 枚举类型 // TODO 未支持多版本
	GenericMap map[int32]interface{}
}

/**
注册契约
*/
func (c *Codec) RegisterContract(code int32, contract interface{}) {
	c.contractMap[code] = contract
}

/**
注册枚举
*/
func (c *Codec) RegisterEnum(code int32, enum interface{}) {
	c.enumMap[code] = enum
}

/**
注册泛型
*/
func (c *Codec) RegisterGeneric(super interface{}, concrete interface{}) {
}

/**
  编码
*/
func (c *Codec) Encode(obj interface{}) ([]byte, error) {
	var err error
	rt := reflect.TypeOf(obj)
	rv := reflect.ValueOf(obj)

	// 编码头信息
	buf := encodeHeader(obj)

	// 编码字段信息
	for i := 0; i < rt.NumField(); i++ {
		tField := rt.Field(i)
		vField := rv.Field(i)
		_, _, _, primitiveType, refContract, refEnum, genericContract, _, numberMask, repeatable, err := resolveTags(tField)
		if err != nil {
			return nil, err
		}

		if primitiveType == PRIMITIVETYPE_BYTES { // 字节数组
			buf = append(buf, encodeBytes(vField.Bytes())...)
		} else {
			repeat := 1
			if repeatable {
				repeat = rv.Field(i).Len()
			}

			if repeat > 1 {
				// 编码数组头信息
				buf = append(buf, encodeArrayHeader(repeat)...)
			}

			for j := 0; j < repeat; j++ {
				var value reflect.Value
				if !repeatable {
					value = vField
				} else {
					value = vField.Index(j)
				}
				if genericContract { // 泛型编码
					buf = append(buf, encodeGeneric(refContract)...)
				} else if refContract != 0 { // 引用其他契约
					buf = append(buf, encodeContract(refContract)...)
				} else if refEnum != 0 { // 引用枚举
					buf = append(buf, encodeEnum(refEnum)...)
				} else { // 基础类型字段
					buf = append(buf, encodePrimitiveType(value, primitiveType, numberMask)...)
				}
			}
		}
	}

	return buf, err
}

/**
解码
*/
func (c *Codec) Decode(data []byte, obj interface{}) {
}

// 解析Tag
func resolveTags(field reflect.StructField) (
	name string,
	order int,
	description string,
	primitiveType string,
	refContract int,
	refEnum int,
	genericContract bool,
	maxSize int,
	numberMask bytes.NumberMask,
	repeatable bool,
	err error) {
	name, ok := field.Tag.Lookup(TAG_NAME)
	if !ok {
		name = ""
	}
	orderStr, ok := field.Tag.Lookup(TAG_ORDER)
	if ok {
		order, err = strconv.Atoi(orderStr)
		if err != nil {
			return
		}
	}
	description, ok = field.Tag.Lookup(TAG_DESCRIPTION)
	if !ok {
		description = ""
	}
	primitiveType = field.Tag.Get(TAG_PRIMITIVETYPE)
	refContractStr, ok := field.Tag.Lookup(TAG_REFCONTRACT)
	if ok {
		refContract, err = strconv.Atoi(refContractStr)
		if err != nil {
			return
		}
	}
	refEnumStr, ok := field.Tag.Lookup(TAG_REFENUM)
	if ok {
		refEnum, err = strconv.Atoi(refEnumStr)
		if err != nil {
			return
		}
	}
	genericContractStr, ok := field.Tag.Lookup(TAG_GENERICCONTRACT)
	if ok && genericContractStr == "true" {
		genericContract = true
	}
	maxSizeStr, ok := field.Tag.Lookup(TAG_MAXSIZE)
	if ok {
		maxSize, err = strconv.Atoi(maxSizeStr)
		if err != nil {
			return
		}
	} else {
		maxSize = -1
	}
	numberEncodingStr, ok := field.Tag.Lookup(TAG_NUMBERENCODING)
	if !ok {
		numberMask = bytes.NUMBERMASK_NONE
	} else {
		numberMask = bytes.GetNumberMask(numberEncodingStr)
	}
	repeatableStr, ok := field.Tag.Lookup(TAG_REPEATABLE)
	if ok && repeatableStr == "true" {
		repeatable = true
	}

	return
}