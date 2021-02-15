package conf

import (
	"gopkg.in/ini.v1"
	"log"
)

var IniFile *ini.File

func init()  {
	cfg, ini_err := ini.Load("conf.ini")
	IniFile = cfg
	if ini_err != nil {
		log.Fatalf("Fail to read file: %v", ini_err)
	}

}