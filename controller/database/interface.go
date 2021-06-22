package database

import (
	"qlite/node"
	class "qlite/struct"
	"viv/static"
)

var Success = []byte("success")

func (m *Manager) Set(name, t []byte) class.Message{
	n := m.list.Work(t, name)
	if n == nil{
		return errStructure
	}
	err := m.table.Set(n)
	if err != nil{
		return class.NewError(err)
	}
	return class.Message{
		Type:    static.ResTypeNode,
		IsWrote: true,
		Node:    n,
	}
}

func (m *Manager) SetX(name, t []byte) class.Message{
	n := m.list.Work(t, name)
	if n == nil{
		return errStructure
	}
	m.table.SetX(n)
	return class.Message{
		Type:    static.ResTypeNode,
		IsWrote: true,
		Node:    n,
	}
}

func (m *Manager) Get(name []byte) class.Message{
	n, err := m.table.Get(name)
	if err != nil{
		return class.NewError(err)
	}
	return class.Message{
		Type:    static.ResTypeNode,
		IsWrote: false,
		Node:    n,
	}
}

func (m *Manager) Type(name []byte) class.Message{
	n, err := m.table.Get(name)
	if err != nil{
		return class.NewError(err)
	}
	return class.Message{
		Type:    static.ResTypeString,
		IsWrote: false,
		String:  n.GetType(),
	}
}

func (m *Manager) Exists(name []byte) class.Message{
	_, err := m.table.Get(name)
	if err != nil{
		return class.Message{
			Type:    static.ResTypeNumber,
			IsWrote: false,
			Number:  0,
		}
	}
	return class.Message{
		Type:    static.ResTypeNumber,
		IsWrote: false,
		Number:  1,
	}
}

func (m *Manager) SetNode(n node.Node) class.Message{
	err := m.table.Set(n)
	if err != nil{
		return class.NewError(err)
	}
	return class.Message{
		Type:    static.ResTypeNode,
		IsWrote: true,
		Node:    n,
	}
}

func (m *Manager) SetNodeX(n node.Node) class.Message{
	m.table.SetX(n)
	return class.Message{
		Type:    static.ResTypeNode,
		IsWrote: true,
		Node:    n,
	}
}

func (m *Manager) Delete(key []byte) class.Message{
	err := m.table.Delete(key)
	if err != nil{
		return class.NewError(err)
	}
	return class.Message{
		Type:    static.ResTypeString,
		IsWrote: true,
		String:  Success,
	}
}

func (m *Manager) DeleteList(key [][]byte) class.Message{
	num := m.table.DeleteList(key)
	return class.Message{
		Type:    static.ResTypeNumber,
		IsWrote: true,
		Number:  num,
	}
}

func (m *Manager) Rename(newKey, oldKey []byte) class.Message{
	err := m.table.Rename(oldKey, newKey)
	if err != nil{
		return class.NewError(err)
	}
	return class.Message{
		Type:    static.ResTypeString,
		IsWrote: true,
		String:  Success,
	}
}

func (m *Manager) RenameX(newKey, oldKey []byte) class.Message{
	err := m.table.RenameX(oldKey, newKey)
	if err != nil{
		return class.NewError(err)
	}
	return class.Message{
		Type:    static.ResTypeString,
		IsWrote: true,
		String:  Success,
	}
}