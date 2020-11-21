package setting_tools

import (
	"fmt"
	"testing"
)

type Setting struct {
	Database *DBSetting `json:"database"`
}

func TestInitSetting(t *testing.T) {
	s := Setting{}
	InitSetting(&s, "./test_setting.json")
	fmt.Println(s.Database.User)
}
