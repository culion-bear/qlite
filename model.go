package main

import "qlite/localTime"

type config struct {
	Password	string	`yaml:"password"`
	IP			string	`yaml:"ip"`
	Port		int		`yaml:"port"`
	Database	int		`yaml:"database"`
	AofPath		string	`yaml:"aof_path"`
	LogPath		string	`yaml:"log_path"`
	StlPath		string	`yaml:"stl_path"`
	AofInterval	int		`yaml:"aof_interval"`
	CORS		bool	`yaml:"cors"`
	TokenKey	string	`yaml:"token_key"`
}

var (
	fileName	=	""
	help		=	false
	version		=	false
	daemon		=	false
	logo 		=
		`
			  ___        _     _ _       
			 / _ \      | |   (_) |_ ___ 
			| | | |_____| |   | | __/ _ \
			| |_| |_____| |___| | ||  __/
			 \__\_\     |_____|_|\__\___|
------------------------------------------------------
		`
)

var lTime = localTime.InitTime()