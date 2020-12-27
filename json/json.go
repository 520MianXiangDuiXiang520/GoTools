package json

import (
	"encoding/json"
	"github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"os"
	"path"
	"reflect"
	"runtime"
)

func load(obj interface{}, path string) {
	if reflect.ValueOf(obj).Elem().Kind() != reflect.Struct {
		panic("obj is not a struct")
	}
	fp, err := os.Open(path)
	if err != nil {
		utils.ExceptionLog(err, "Fail to open obj")
		panic(err)
	}
	defer fp.Close()
	decoder := json.NewDecoder(fp)
	err = decoder.Decode(&obj)
	if err != nil {
		utils.ExceptionLog(err, "Fail to decode json obj")
		panic(err)
	}
}

// FromFileLoadToObj 从 fPath 路径下的 json 文件中加载数据到 s 对象
func FromFileLoadToObj(s interface{}, fPath string) {
	if reflect.ValueOf(s).Kind() != reflect.Ptr {
		panic("The s is not a ptr")
	}
	_, currently, _, _ := runtime.Caller(1)
	filename := path.Join(path.Dir(currently), fPath)
	load(s, filename)
}
