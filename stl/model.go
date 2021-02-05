package api

import "errors"

var(
	ErrServiceExist = errors.New("service's name is exists")
	ErrNotExist     = errors.New("service is not exists")
	ErrApiNotExit	= errors.New("api is not exists")
	ErrApiType		= errors.New("api type is not compare")
	ErrServiceType	= errors.New("service type is not compare")
	ErrIDEmpty		= errors.New("id is empty")
)

const(
	STL	= iota
	ALG
)

type stlServiceInfo struct {
	url			string
	handle		StlClient
	isOrderly	bool
	apiMap		map[string]*ApiInfo
}

var servers map[string]*stlServiceInfo = make(map[string]*stlServiceInfo)