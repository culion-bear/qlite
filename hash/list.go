package hash

//查询key是否过期
func isOverdue(t NodeTime) bool{
	if t.durationTime <= 0 {
		return false
	}
	return lTime.GetTime() >= t.durationTime+t.beginTime
}

//删除一个node
func (handle *List) del(node *nodeList){
	if node.last == nil && node.next == nil{
		handle.head = nil
		handle.tail = nil
	}else if node.last == nil{
		handle.head = node.next
	}else if node.next == nil{
		handle.tail = node.last
	}else{
		node.last.next = node.next
		node.next.last = node.last
	}
}

//查找value
func (handle *List) get(key string) *nodeList{
	handle.lock.Lock()
	defer handle.lock.Unlock()
	for i := handle.head ; i != nil ; i=i.next{
		if isOverdue(i.value.GetTime()) {
			handle.del(i)
			continue
		}
		if i.value.GetKey() == key{
			return i
		}
	}
	return nil
}

//在最后进行添加
func (handle *List) pushBack(node *nodeList){
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.head == nil{
		handle.head = node
		handle.tail = node
	}else{
		handle.tail.next = node
		node.last = handle.tail
	}
}

//添加数据，若存在则返回false
func (handle *List) insert(node Node) error{
	v := handle.get(node.GetKey())
	if v != nil{
		return ErrIsExists
	}
	handle.pushBack(&nodeList{
		last:nil,
		next:nil,
		value:node,
	})
	return nil
}

//添加数据，若存在则覆盖
func (handle *List) insertX(node Node){
	v := handle.get(node.GetKey())
	if v != nil{
		handle.lock.Lock()
		v.value=node
		handle.lock.Unlock()
	}else{
		handle.pushBack(&nodeList{
			last:nil,
			next:nil,
			value:node,
		})
	}
}

//查询所有的key值
func (handle *List) getAll() []string{
	keys := make([]string,0)
	handle.lock.Lock()
	defer handle.lock.Unlock()
	for i := handle.head ; i != nil ; i=i.next{
		if isOverdue(i.value.GetTime()) {
			handle.del(i)
			continue
		}
		keys = append(keys, i.value.GetKey())
	}
	return keys
}

//查询所有key的具体信息
func (handle *List) getAllX() []Info{
	keys := make([]Info,0)
	handle.lock.Lock()
	defer handle.lock.Unlock()
	for i := handle.head ; i != nil ; i=i.next{
		if isOverdue(i.value.GetTime()) {
			handle.del(i)
			continue
		}
		keys = append(keys, Info{
			Key:  i.value.GetKey(),
			Type: i.value.GetType(),
			Time: func(t NodeTime) int64{
				if t.durationTime<=0 {
					return 0
				}
				return t.durationTime+t.beginTime
			}(i.value.GetTime()),
		})
	}
	return keys
}

//获取个数
func (handle *List) size() int{
	sum := 0
	handle.lock.Lock()
	defer handle.lock.Unlock()
	for i := handle.head ; i != nil ; i=i.next{
		if isOverdue(i.value.GetTime()) {
			handle.del(i)
			continue
		}
		sum++
	}
	return sum
}

//获取类别
func (handle *List) getType(key string) (string,error){
	node := handle.get(key)
	if node == nil{
		return "",ErrNotFound
	}
	return node.value.GetType(),nil
}

//是否存在
func (handle *List) isExists(key string) bool{
	return handle.get(key) != nil
}

//设置到期时间
func (handle *List) pex(key string,t int64) error{
	node := handle.get(key)
	if node == nil{
		return ErrNotFound
	}
	node.value.SetTime(NodeTime{
		beginTime:lTime.GetTime(),
		durationTime:t,
	})
	return nil
}

//设置何时到期
func (handle *List) pexTo(key string,t int64) error{
	node := handle.get(key)
	if node == nil{
		return ErrNotFound
	}
	b := lTime.GetTime()
	node.value.SetTime(NodeTime{
		beginTime:b,
		durationTime:t+b,
	})
	return nil
}

//查看剩余时间
func (handle *List) rTime(key string) (int64,error){
	node := handle.get(key)
	if node == nil{
		return 0,ErrNotFound
	}
	t := node.value.GetTime()
	if t.durationTime <= 0{
		return 0,nil
	}
	return t.durationTime+t.beginTime-lTime.GetTime(),nil
}

//查看结束时间
func (handle *List) eTime(key string) (int64,error){
	node := handle.get(key)
	if node == nil{
		return 0,ErrNotFound
	}
	t := node.value.GetTime()
	if t.durationTime <= 0{
		return 0,nil
	}
	return t.durationTime+t.beginTime,nil
}

//删除
func (handle *List) delete(key string) error{
	node := handle.get(key)
	if node == nil{
		return ErrNotFound
	}
	handle.lock.Lock()
	defer handle.lock.Unlock()
	handle.del(node)
	return nil
}