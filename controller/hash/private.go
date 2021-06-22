package hash

import class "qlite/struct"

func (m *Manager) push() error{
	var err error
	err = m.list.Push([]byte("set"), m.set)
	if err != nil{
		return err
	}
	err = m.list.Push([]byte("setx"), m.setX)
	if err != nil{
		return err
	}
	err = m.list.Push([]byte("del"), m.del)
	if err != nil{
		return err
	}
	err = m.list.Push([]byte("delx"), m.delList)
	if err != nil{
		return err
	}
	err = m.list.Push([]byte("rename"), m.rename)
	if err != nil{
		return err
	}
	err = m.list.Push([]byte("renamex"), m.renameX)
	if err != nil{
		return err
	}
	err = m.list.Push([]byte("type"), m.nodeType)
	if err != nil{
		return err
	}
	err = m.list.Push([]byte("exists"), m.exists)
	if err != nil{
		return err
	}
	return m.list.Push([]byte("get"), m.get)
}

func (m *Manager) set(msg []class.Message) class.Message{
	if len(msg) != 2{
		return errPackage
	}
	return m.handle.Set(msg[0].ToString(), msg[1].ToString())
}

func (m *Manager) setX(msg []class.Message) class.Message{
	if len(msg) != 2{
		return errPackage
	}
	return m.handle.SetX(msg[0].ToString(), msg[1].ToString())
}

func (m *Manager) del(msg []class.Message) class.Message{
	if len(msg) != 1{
		return errPackage
	}
	return m.handle.Delete(msg[0].ToString())
}

func (m *Manager) delList(msg []class.Message) class.Message{
	if len(msg) != 1{
		return errPackage
	}
	if !msg[0].IsList(){
		return errPackage
	}
	return m.handle.DeleteList(msg[0].List)
}

func (m *Manager) rename(msg []class.Message) class.Message{
	if len(msg) != 2{
		return errPackage
	}
	return m.handle.Rename(msg[0].ToString(), msg[1].ToString())
}

func (m *Manager) renameX(msg []class.Message) class.Message{
	if len(msg) != 2{
		return errPackage
	}
	return m.handle.RenameX(msg[0].ToString(), msg[1].ToString())
}

func (m *Manager) get(msg []class.Message) class.Message{
	if len(msg) != 1{
		return errPackage
	}
	return m.handle.Get(msg[0].ToString())
}

func (m *Manager) nodeType(msg []class.Message) class.Message{
	if len(msg) != 1{
		return errPackage
	}
	return m.handle.Type(msg[0].ToString())
}

func (m *Manager) exists(msg []class.Message) class.Message{
	if len(msg) != 1{
		return errPackage
	}
	return m.handle.Exists(msg[0].ToString())
}