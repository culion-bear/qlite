package hash

import (
	"errors"
	"qlite/localTime"
	"sync"
)

var lTime = localTime.InitTime()
var rootDataBase []*Hash

var(
	ErrNotFound = errors.New("key is not found")
	ErrIsExists = errors.New("key is exists")
	ErrDataBase = errors.New("database num is bad")
	ErrNotHash  = errors.New("key's type is not hash")
)

type Node interface {
	GetKey() string
	SetKey(string)
	GetTime() NodeTime
	SetTime(NodeTime)
	GetType() string
	GetID() string
	UpdateID(string)
}

type Info struct {
	Key				string	`json:"key"`
	Type			string	`json:"type"`
	Time			int64	`json:"time"`
}

type Hash struct {
	key				string
	lock			sync.RWMutex
	time			NodeTime
	value			nodeHash
}

type Other struct {
	nodeType		string
	key				string
	id				string
	lock			sync.RWMutex
	time			NodeTime
}

type List struct {
	lock			sync.Mutex
	head			*nodeList
	tail			*nodeList
}

type nodeHash struct {
	list			[65536]*List
}

type nodeList struct {
	last			*nodeList
	next			*nodeList
	value			Node
}

type NodeTime struct {
	durationTime	int64
	beginTime		int64
}

func (handle *Hash) GetKey() string{
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	return handle.key
}

func (handle *Hash) SetKey(key string){
	handle.lock.Lock()
	defer handle.lock.Unlock()
	handle.key=key
}

func (handle *Hash) GetTime() NodeTime{
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	return handle.time
}

func (handle *Hash) SetTime(time NodeTime){
	handle.lock.Lock()
	defer handle.lock.Unlock()
	handle.time=time
}

func (handle *Hash) GetType() string{
	return "hash"
}

func (handle *Hash) GetID() string{
	return ""
}

func (handle *Hash) UpdateID(_ string){}

func (handle *Other) GetKey() string{
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	return handle.key
}

func (handle *Other) SetKey(key string){
	handle.lock.Lock()
	defer handle.lock.Unlock()
	handle.key=key
}

func (handle *Other) GetTime() NodeTime{
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	return handle.time
}

func (handle *Other) SetTime(time NodeTime){
	handle.lock.Lock()
	defer handle.lock.Unlock()
	handle.time=time
}

func (handle *Other) GetType() string{
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	return handle.nodeType
}

func (handle *Other) GetID() string{
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	return handle.id
}

func (handle *Other) UpdateID(id string){
	handle.lock.Lock()
	defer handle.lock.Unlock()
	handle.id=id
}