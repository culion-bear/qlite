package static

import (
	"errors"
	class "qlite/struct"
)

type Error []byte

func (m Error) String() string{
	return string(m)
}

var(
	//hash error
	ErrNodeIsExist		=	Error("node is exist")
	ErrNodeNotExist		=	Error("node is not exist")
)

var(
	//scanner error
	ErrNumber			=	Error("package length is error")
)

var(
	//new node error
	ErrStructure		=	Error("structure is not exist")
)

var(
	//password error
	ErrNotSplitPoint	=	errors.New("split point is not found")
	ErrPasswordError	=	errors.New("password is not true")
)

var(
	//resolver error
	ErrTypeNotExist		=	Error("type is not exist")
	ErrNotResPoint		=	Error("split point is not found")
	ErrPackage			=	Error("package is error")
)

var(
	//scheduler error
	ErrWithOutReturn	=	Error("!a;not return")
	ErrNotComplete		=	class.Message{
								Type:    '!',
								String:  Error("package is not complete"),
							}
)