package setting_tools

import (
	"encoding/json"
	"github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"os"
	"path"
	"reflect"
	"runtime"
)

func load(setting interface{}, path string) {
	if reflect.ValueOf(setting).Elem().Kind() != reflect.Struct {
		panic("setting is not a struct")
	}
	fp, err := os.Open(path)
	if err != nil {
		utils.ExceptionLog(err, "Fail to open setting")
		panic(err)
	}
	defer fp.Close()
	decoder := json.NewDecoder(fp)
	err = decoder.Decode(&setting)
	if err != nil {
		utils.ExceptionLog(err, "Fail to decode json setting")
		panic(err)
	}
}

func InitSetting(s interface{}, fPath string) {
	if reflect.ValueOf(s).Kind() != reflect.Ptr {
		panic("The s is not a ptr")
	}
	_, currently, _, _ := runtime.Caller(1)
	filename := path.Join(path.Dir(currently), fPath)
	load(s, filename)
}
