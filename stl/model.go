package api

import (
	"errors"
	"qlite/localTime"
	"sync"
)

var(
	ErrServiceExist = errors.New("service's name is exists")
	ErrNotExist     = errors.New("service is not exists")
	ErrApiNotExit	= errors.New("api is not exists")
	ErrApiType		= errors.New("api type is not compare")
	ErrServiceType	= errors.New("service type is not compare")
	ErrIDEmpty		= errors.New("id is empty")
	ErrServiceToken	= errors.New("password error")
	ErrRestore		= errors.New("service is restoring")
)

const(
	STL	= iota
	ALG
)

type stlServiceInfo struct {
	url			string
	name		string
	handle		StlClient
	isOrderly	bool
	apiMap		map[string]*ApiInfo
	token		string
	flag		bool
	mu			sync.RWMutex
}

var servers = make(map[string]*stlServiceInfo)

var AofPath string

var lTime =localTime.InitTime()