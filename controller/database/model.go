package database

import (
	class "qlite/struct"
	"viv/static"
)

var errStructure = class.Message{
	Type:    static.ResTypeError,
	IsWrote: false,
	String:  static.ErrStructure,
}