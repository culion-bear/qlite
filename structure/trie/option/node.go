package option

import (
	"errors"
	class "qlite/struct"
	"qlite/trie"
)

type node struct {
	list	[26]*node
	value	trie.TrieManager
}

var errOptionNotFound = class.Message{
	Type:    '!',
	IsWrote: false,
	String:  []byte("structure is not found"),
}


var errOptionIllegal = class.Message{
	Type:    '!',
	IsWrote: false,
	String:  []byte("structure is illegal"),
}

func newNode() *node{
	return &node{}
}

func (m *node) push(msg []byte, f trie.TrieManager, k int) error{
	if k == len(msg){
		if m.value != nil{
			return errors.New(string(msg) + " structure is exist")
		}
		m.value = f
		return nil
	}
	key := msg[k] - 'a'
	if key >= 26 || key < 0{
		return errors.New(string(msg) + " structure is illegal")
	}
	if m.list[key] == nil{
		m.list[key] = newNode()
	}
	return m.list[key].push(msg, f, k + 1)
}

func (m *node) work(name, opt []byte, msg []class.Message, k, length int) class.Message{
	if k == length{
		if m.value == nil{
			return errOptionNotFound
		}
		return m.value.Work(opt, msg)
	}
	key := name[k] - 'a'
	if key >= 26 || key < 0{
		return errOptionIllegal
	}
	if m.list[key] == nil {
		return errOptionNotFound
	}
	return m.list[key].work(name, opt, msg, k + 1, length)
}