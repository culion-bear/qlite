package config

import "encoding/xml"

type Config struct {
	XMLName		xml.Name	`xml:"config"`
	Password	string		`xml:"password"`
	IP			string		`xml:"ip"`
	Port		int			`xml:"port"`
	Aof			string		`xml:"aof"`
	Plugins		string		`xml:"plugins"`
	Pid			string		`xml:"pid"`
	Interval	int			`xml:"interval"`
}