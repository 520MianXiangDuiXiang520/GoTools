package json

import (
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
	FromFileLoadToObj(&s, "./test_setting.json")
	if s.Database.Port != 3306 {
		t.Error("db")
	}
}
