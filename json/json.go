package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"github.com/520MianXiangDuiXiang520/GinTools/path_tools"
	"os"
	"path"
	"reflect"
	"runtime"
)

func load(obj interface{}, path string) {
	fp, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("Fail to open %s", path)
		utils.ExceptionLog(err, msg)
		panic(err)
	}
	defer fp.Close()
	decoder := json.NewDecoder(fp)
	err = decoder.Decode(&obj)
	if err != nil {
		msg := fmt.Sprintf("Fail to decode data to %v", obj)
		utils.ExceptionLog(err, msg)
		panic(err)
	}
}

// FromFileLoadToObj 从 fPath 路径下的 json 文件中加载数据到 s 对象
func FromFileLoadToObj(s interface{}, fPath string) {
	if reflect.ValueOf(s).Kind() != reflect.Ptr {
		msg := "The s must be a ptr type"
		utils.ExceptionLog(errors.New("TypeError"), msg)
		panic("The s is not a ptr")
	}
	if !path_tools.IsAbs(fPath) {
		_, currently, _, _ := runtime.Caller(1)
		fPath = path.Join(path.Dir(currently), fPath)
	}
	load(s, fPath)
}
