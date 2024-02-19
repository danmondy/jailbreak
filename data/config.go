package data

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host         string       `json:"host"`
	Port         string       `json:"port"`
	Db           DbConnection `json:"db"`
	TemplatePath string       `json:"templateFolder"`
	PublicDir    string       `json:"publicFolder"`
}

func ReadConfig(file string, c interface{}) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, c)
}

type DbConnection struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
