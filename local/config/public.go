package config

import (
	"encoding/xml"
	"viv/local/file"
)

func Read(path string) (Config, error){
	msg, err := file.ReadFile(path)
	var c Config
	if err != nil{
		return c, err
	}
	err = xml.Unmarshal(msg, &c)
	return c, err
}