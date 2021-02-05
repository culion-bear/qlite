package hash

func newList () *List{
	return &List{
		head:nil,
		tail:nil,
	}
}

func (handle *Hash) Get(key string) (Node,error){
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return nil,ErrNotFound
	}
	n := handle.value.list[k].get(key)
	if n == nil{
		return nil,ErrNotFound
	}
	return n.value,nil
}

func (handle *Hash) Set(node Node) error{
	k := toHash(node.GetKey())
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.value.list[k] == nil{
		handle.value.list[k] = newList()
	}
	return handle.value.list[k].insert(node)
}

func (handle *Hash) SetX(node Node){
	k := toHash(node.GetKey())
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.value.list[k] == nil{
		handle.value.list[k] = newList()
	}
	handle.value.list[k].insertX(node)
}

func (handle *Hash) SelectKeyName() []string{
	names := make([]string,0)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	for _,value := range handle.value.list{
		if value == nil{
			continue
		}
		names = append(names,value.getAll()...)
	}
	return names
}

func (handle *Hash) SelectInfo() []Info{
	infos := make([]Info,0)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	for _,value := range handle.value.list{
		if value == nil{
			continue
		}
		infos = append(infos,value.getAllX()...)
	}
	return infos
}

func (handle *Hash) Del(keys []string) int{
	sum := 0
	for _,v := range keys{
		k := toHash(v)
		handle.lock.RLock()
		if handle.value.list[k] == nil{
			handle.lock.RUnlock()
			continue
		}
		handle.lock.RUnlock()
		if handle.value.list[k].delete(v) == nil{
			sum++
		}
	}
	return sum
}

func (handle *Hash) GetNodeType(key string) (string,error){
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return "",ErrNotFound
	}
	return handle.value.list[k].getType(key)
}

func (handle *Hash) GetNodeID(key string) (string,error){
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return "",ErrNotFound
	}
	n := handle.value.list[k].get(key)
	if n == nil{
		return "",ErrNotFound
	}
	return n.value.GetID(),nil
}

func (handle *Hash) UpdateNodeID(key,id string) error{
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return ErrNotFound
	}
	handle.value.list[k].get(key).value.UpdateID(id)
	return nil
}

func (handle *Hash) Exists(key string) bool{
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return false
	}
	return handle.value.list[k].isExists(key)
}

func (handle *Hash) Pex(key string,t int64) error{
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return ErrNotFound
	}
	return handle.value.list[k].pex(key,t)
}

func (handle *Hash) PexTo(key string,t int64) error{
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return ErrNotFound
	}
	return handle.value.list[k].pexTo(key,t)
}

func (handle *Hash) ETime(key string) (int64,error){
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return 0,ErrNotFound
	}
	return handle.value.list[k].eTime(key)
}

func (handle *Hash) RTime(key string) (int64,error){
	k := toHash(key)
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	if handle.value.list[k] == nil{
		return 0,ErrNotFound
	}
	return handle.value.list[k].rTime(key)
}

func (handle *Hash) Rename(key,newKey string) error{
	k,nK := toHash(key),toHash(newKey)
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.value.list[k] == nil{
		return ErrNotFound
	}
	if handle.value.list[nK] != nil && handle.value.list[nK].isExists(newKey){
		return ErrIsExists
	}
	nodeKey := handle.value.list[k].get(key)
	if nodeKey == nil{
		return ErrNotFound
	}
	n := nodeKey.value
	handle.value.list[k].del(nodeKey)
	n.SetKey(newKey)
	if handle.value.list[nK] == nil{
		handle.value.list[nK] = newList()
	}
	return handle.value.list[nK].insert(n)
}

func (handle *Hash) RenameX(key,newKey string) error{
	k,nK := toHash(key),toHash(newKey)
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.value.list[k] == nil{
		return ErrNotFound
	}
	nodeKey := handle.value.list[k].get(key)
	if nodeKey == nil{
		return ErrNotFound
	}
	n := nodeKey.value
	handle.value.list[k].del(nodeKey)
	n.SetKey(newKey)
	if handle.value.list[nK] == nil{
		handle.value.list[nK] = newList()
	}
	handle.value.list[nK].insertX(n)
	return nil
}

func (handle *Hash) Size() int{
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	sum := 0
	for _,v := range handle.value.list{
		if v == nil{
			continue
		}
		sum += v.size()
	}
	return sum
}