package network

import (
	"errors"
	"github.com/kataras/iris/v12"
	"qlite/localTime"
	"qlite/persistence"
)

const(
	ErrJson    = 400
	ErrService = 500
	ErrToken   = 403
	ErrLogin   = 401
	ErrPing    = 202
	Success    = 200
)

const Version = "V 2.1.2"

var(
	AofHandle *persistence.AofManager
	LogHandle *persistence.LogManager
	StlHandle *persistence.StlManager
)

var(
	ErrKeyEmpty = errors.New("key is empty")
)

var lTime = localTime.InitTime()

var irisJson = iris.JSONReader{
	DisallowUnknownFields: true,
}

var TokenKey string

var Password string

type User struct {
	Password	string		`json:"password"`
}

type Error struct {
	Code		int			`json:"code"`
	Msg			string		`json:"msg"`
}

type StlUrl struct {
	Url			string		`json:"url"`
	Password	string		`json:"password,omitempty"`
}

//type SetModel struct {
//	Key			string	`json:"key"`
//	Type		string	`json:"type"`
//	Time		int64	`json:"time"`
//}

type KeyList struct {
	Keys		[]string	`json:"keys"`
}

type KeyModel struct {
	Key			string		`json:"key"`
}

type TimeModel struct {
	Key			string		`json:"key"`
	Time		int64		`json:"time"`
}

type RenameModel struct {
	Key			string		`json:"key"`
	NewKey		string		`json:"new_key"`
}

type StlApiModel struct {
	Key			string		`json:"key"`
	Time		int64		`json:"time,omitempty"`
	Opt			interface{}	`json:"opt,omitempty"`
}

type KeyAddress struct {
	Key			string		`json:"key"`
	Path		string		`json:"path,omitempty"`
}

type AlgApiModel struct {
	Keys		[]KeyAddress`json:"keys,omitempty"`
	Opt			interface{}	`json:"opt,omitempty"`
}