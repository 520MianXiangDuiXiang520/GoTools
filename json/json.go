package json

import (
	"encoding/json"
	"fmt"
	path2 "github.com/520MianXiangDuiXiang520/GoTools/path"
	"os"
	"path"
	"reflect"
	"runtime"
)

func load(obj interface{}, path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("fail to open %s, error: %w", path, err)
	}
	defer fp.Close()
	decoder := json.NewDecoder(fp)
	err = decoder.Decode(&obj)
	if err != nil {
		return fmt.Errorf("fail to decode data to %v, error: %w", obj, err)
	}
	return nil
}

// FromFileLoadToObj 从 fPath 路径下的 json 文件中加载数据到 s 对象
func FromFileLoadToObj(s interface{}, fPath string) error {
	if reflect.ValueOf(s).Kind() != reflect.Ptr {
		return fmt.Errorf("the %v is not a ptr", s)
	}
	if !path2.IsAbs(fPath) {
		_, currently, _, _ := runtime.Caller(1)
		fPath = path.Join(path.Dir(currently), fPath)
	}
	return load(s, fPath)
}
