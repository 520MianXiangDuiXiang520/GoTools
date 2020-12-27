package setting_tools

import (
	"fmt"
	"testing"
)

type DBSetting struct {
	Engine   string `json:"engine"`
	DBName   string `json:"db_name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type Setting struct {
	Database *DBSetting `json:"database"`
}

func TestInitSetting(t *testing.T) {
	s := Setting{}
	InitSetting(&s, "./test_setting.json")
	fmt.Println(s.Database.User)
}
