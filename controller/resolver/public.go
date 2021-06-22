package resolver

import (
	"viv/static"
	"viv/structure/qt"
)

func (m *Manager) Split(buf []byte) (*qt.Queue, static.Error){
	return m.split(buf)
}