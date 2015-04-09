package funcMap

import (
	"reflect"
	"unicode"
)

const (
	ERROR_CODE_ARGS_NUM = 1000000 + iota
	ERROR_CODE_ARGS_TYPE
)

type FuncMap struct {
	funcParamsMap map[string][]string    //方法参数
	funcOwnerMap  map[string]interface{} //方法拥有
}

func NewFuncMap() *FuncMap {
	return &FuncMap{
		funcParamsMap: make(map[string][]string),
		funcOwnerMap:  make(map[string]interface{}),
	}
}

//大写开头为可导出方法
func IsExportedName(name string) bool {
	return name != "" && unicode.IsUpper(rune(name[0]))
}

//获得可导出方法的方法名及参数类型
func GetMethods(v interface{}) map[string][]string {
	funcMap := make(map[string][]string)
	reflectType := reflect.TypeOf(v)
	for i := 0; i < reflectType.NumMethod(); i++ {
		funcList := make([]string, 0)
		method := reflectType.Method(i)
		methodType := method.Type
		methodName := method.Name
		if !IsExportedName(methodName) {
			continue
		}
		for j := 1; j < methodType.NumIn(); j++ {
			params := methodType.In(j)
			funcList = append(funcList, params.String())
		}
		funcMap[methodName] = funcList
	}
	return funcMap
}

//执行...
func (f FuncMap) Invoke(funcName string, args ...interface{}) (result []interface{}, err interface{}) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = err1
		}
	}()
	if receiver, ok := f.funcOwnerMap[funcName]; ok {
		if pts, ok := f.funcParamsMap[funcName]; ok {
			//检查参数长度
			if len(pts) != len(args) {
				err = ERROR_CODE_ARGS_NUM //参数数量错误
				return
			}
			//检查参数类型
			for index, paramType := range pts {
				if reflect.TypeOf(args[index]).String() != paramType {
					err = ERROR_CODE_ARGS_TYPE //参数类型错误
					return
				}
			}

			inputs := make([]reflect.Value, len(args))
			for i, _ := range args {
				inputs[i] = reflect.ValueOf(args[i])
			}
			rv := reflect.ValueOf(receiver).MethodByName(funcName).Call(inputs)
			result = make([]interface{}, len(rv))
			for k, v := range rv {
				result[k] = v.Interface()
			}
		}
	}

	return
}

func (f FuncMap) Register(v interface{}) {
	fs := GetMethods(v)
	for key, value := range fs {
		f.funcParamsMap[key] = value
		f.funcOwnerMap[key] = v
	}
}
