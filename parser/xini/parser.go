package xini

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// @Date：   2018/8/19 0019 10:54
// @Author:  Joshua Conero
// @Name:    解析器

// Parser the ini file base parse interface
type Parser interface {
	Get(key string) (bool, any)
	GetDef(key string, def any) any
	HasKey(key string) bool
	// SetFunc 函数式值获取
	// 实现如动态值获取，类似 js 对象中的 [get function()]
	SetFunc(key string, regFn func() any) Parser

	// Raw 支持多级数据访问，获取元素数据
	// 实际读取的原始数据为 map[string]string
	Raw(key string) string

	// Value get or set value: key, value(nil), default
	Value(params ...any) any

	GetAllSection() []string
	// Section the param format support
	//		1.     fun Section(section, key string) 	二级访问
	//		2.     fun Section(format string) 			点操作
	Section(params ...any) any

	GetData() map[string]any

	Set(key string, value any) Parser // 设置值
	Del(key string) bool              // 删除键值

	IsValid() bool
	OpenFile(filename string) (error, []ScannerLog)
	ReadStr(content string) Parser
	ErrorMsg() string // 错误信息

	Save() bool
	SaveAsFile(filename string) bool
	Driver() string
}

// 解析为数字类型，i64/f64
func parseNumber(vStr string) (value any, isOk bool) {
	i64Symbol := getRegByKey("reg_i64_symbol")
	if i64Symbol != nil && i64Symbol.MatchString(vStr) {
		i64, er := strconv.ParseInt(vStr, 10, 64)
		if er == nil {
			value = i64
			isOk = true
			return
		}
	}

	f64Symbol := getRegByKey("reg_f64_symbol")
	if f64Symbol != nil && f64Symbol.MatchString(vStr) {
		f64, er := strconv.ParseFloat(vStr, 64)
		if er == nil {
			value = f64
			isOk = true
			return
		}
	}
	return
}

// 字符串清理
func stringClear(vStr string) string {
	vStr = strings.TrimSpace(vStr)
	strSymbol := getRegByKey("reg_str_symbol")
	if strSymbol != nil && strSymbol.MatchString(vStr) {
		vStr = vStr[1 : len(vStr)-1]
	}
	return vStr
}

// 行内字符清理，如含内注释等
func lnTrim(vStr string) (vs string) {
	vs = vStr
	vs = strings.TrimSpace(vs)

	strSymbol := getRegByKey("reg_str_symbol")
	// 标准的`"string"`或`'string'`不再进行处理
	if strSymbol != nil && strSymbol.MatchString(vs) {
		return
	}

	lnReg := getRegByKey("reg_str_symbol_ln")
	if lnReg != nil && lnReg.MatchString(vs) {
		line := lnReg.FindAllString(vs, -1)
		dick := map[string]string{}
		tmStr := fmt.Sprintf("%v", time.Now().Unix())
		for i, ln := range line {
			key := fmt.Sprintf("L%vN%v", tmStr, i)
			vs = strings.ReplaceAll(vs, ln, key)
			dick[key] = ln
		}

		cmtReg := getRegByKey("reg_has_comment")
		if cmtReg != nil && cmtReg.MatchString(vs) {
			indexList := cmtReg.FindAllStringIndex(vs, -1)
			if len(indexList) > 0 {
				if len(indexList[0]) > 0 {
					vs = vs[:indexList[0][0]]
				}
			}
		}

		// 字符串还原
		for rpl, raw := range dick {
			vs = strings.ReplaceAll(vs, rpl, raw)
		}
		vs = strings.TrimSpace(vs)
	}

	return
}

// 切片解析（行内）
func parseSlice(vStr string) (value any, isOk bool) {
	if !strings.Contains(vStr, baseLimiterToken) {
		return
	}

	dick := map[string]string{}
	lnReg := getRegByKey("reg_str_symbol_ln")
	isStr := false
	if lnReg != nil && lnReg.MatchString(vStr) {
		line := lnReg.FindAllString(vStr, -1)
		tmStr := fmt.Sprintf("%v", time.Now().Unix())
		for i, ln := range line {
			key := fmt.Sprintf("L%vN%v", tmStr, i)
			vStr = strings.ReplaceAll(vStr, ln, key)
			dick[key] = ln
		}
		isStr = true
	}

	// 字符换分隔
	var strQue []string
	var i64Que []int64
	var f64Que []float64
	var isI64 = false
	for _, s := range strings.Split(vStr, baseLimiterToken) {
		// 字符串
		if isStr {
			for rpl, raw := range dick {
				s = strings.ReplaceAll(s, rpl, raw)
			}
			s = stringClear(s)
			strQue = append(strQue, s)
		} else {
			num, isNum := parseNumber(strings.TrimSpace(s))
			if isNum {
				if i64, isMatch := num.(int64); isMatch {
					if !isI64 && len(i64Que) == 0 && len(f64Que) == 0 {
						isI64 = true
					}
					if isI64 {
						i64Que = append(i64Que, i64)
					} else {
						// i64 存在且不是i64是将原数组切换到f64存储栈中
						for _, oldI64 := range i64Que {
							f64Que = append(f64Que, float64(oldI64))
						}
						i64Que = []int64{}
					}
				} else {
					f64Que = append(f64Que, num.(float64))
					isI64 = false
				}
			} else {
				s = stringClear(s)
				strQue = append(strQue, s)
			}

		}
	}

	if len(strQue) > 0 {
		isOk = true
		value = strQue
	} else if len(f64Que) > 0 {
		isOk = true
		value = f64Que
	} else if len(i64Que) > 0 {
		isOk = true
		value = i64Que
	}

	return
}

// 将字符串解析为参数
// 将原始的字符串解析为对应的参数
func parseValue(vStr string) any {
	var value any
	switch strings.ToLower(vStr) {
	case "true":
		value = true
	case "false":
		value = false
	default:
		// 包裹找字符串如，`"string"` 或 `'string'`
		if v, isOk := parseNumber(vStr); isOk {
			value = v
		} else if v, isOk = parseSlice(vStr); isOk {
			value = v
		} else {
			value = stringClear(vStr)
		}
	}
	return value
}

// DecKvPairs Deconstruct key value pairs
func DecKvPairs(s string) *KvPairs {
	var kp = &KvPairs{
		raw: s,
	}
	ln := strings.TrimSpace(s)
	ln = lnTrim(ln) // 行内注释处理
	strSymbol := getRegByKey("reg_str_symbol")
	if strSymbol != nil && strSymbol.MatchString(ln) {
		kp.isString = true
		kp.value = s[1 : len(s)-1]
	} else {
		eqlIdx := strings.Index(ln, baseEqualToken)
		if eqlIdx > -1 {
			kp.isKv = true
			kp.key = strings.TrimSpace(ln[:eqlIdx])
			kp.value = strings.TrimSpace(ln[eqlIdx+1:])
		} else {
			kp.value = strings.TrimSpace(ln)
			kp.isString = true
		}
	}
	return kp
}

// KvPairs KV 键值对
type KvPairs struct {
	key      string
	value    string
	isString bool
	isKv     bool
	raw      string
}
