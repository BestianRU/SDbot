package cfg

import (
	"encoding/json"
	"os"
)

//TConnection Telegram connection structure
type TConnection struct {
	Token   string `json:"token"`
	Timeout int    `json:"timeout"`
	Debug   bool   `json:"debug"`
}

//MsqlConnection mysql connection structure
type MsqlConnection struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
}

//Cfg structure for sdbot
type Cfg struct {
	T          TConnection    `json:"telegram"`
	M          MsqlConnection `json:"mysql"`
	AuthUser   string         `json:"authUser"`
	MsgNotAuth string         `json:"msgNotAuth"`
}

//Load config from "./sdbotcfg.json"
func (c *Cfg) Load() error {
	file, err := os.Open("sdbotcfg.json")
	defer file.Close()
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(c)
	if err != nil {
		return err
	}
	return nil
}
