package controller

import (
	"errors"
	"plugin"
	"qlite/node"
	"qlite/trie"
	"viv/controller/database"
	"viv/controller/hash"
	"viv/controller/logistics"
	sd "viv/controller/scheduler"
	"viv/local/config"
	"viv/local/file"
)

func Init(c config.Config) error{
	sd.Scheduler = sd.New([]byte(c.Password))
	list, err := initDB(c.Plugins)
	if err != nil{
		return err
	}
	for _, v := range list{
		err = sd.Scheduler.Push(v)
		if err != nil{
			return err
		}
	}
	log, err := logistics.New(c.Aof, c.Interval, sd.Scheduler.Load)
	if err != nil{
		return nil
	}
	err = log.Load()
	if err != nil{
		return err
	}
	go sd.Scheduler.Run(log)
	return nil
}

func initDB(path string) ([]trie.TrieManager, error){
	db, list := database.New(), make([]trie.TrieManager, 0)
	err := db.Push([]byte("object"), node.NewObject)
	if err != nil{
		return list, err
	}

	n, err := hash.New(db)
	if err != nil{
		return list, err
	}
	list = append(list, n)

	l, err := initPlugins(path, db)
	if err != nil{
		return list, err
	}

	return append(list, l...), nil
}

func initPlugins(path string, db *database.Manager) ([]trie.TrieManager, error){
	names, err := file.GetFileName(path)
	if err != nil{
		return nil, err
	}
	l := make([]trie.TrieManager, 0)
	for _, v := range names{
		t, err := load(v, db)
		if err != nil{
			return nil, err
		}
		l = append(l, t)
	}
	return l, nil
}

func load(path string, db *database.Manager) (trie.TrieManager, error){
	model, err := plugin.Open(path)
	if err != nil{
		return nil, err
	}
	f, err := getNewTireManagerFunc(model, path)
	if err != nil{
		return nil, err
	}
	t, err := f(db)
	if err != nil{
		return nil, err
	}
	n := getNewNodeFunc(model)
	if n != nil{
		err = db.Push(t.GetName(), n)
		if err != nil{
			return nil, err
		}
	}
	return t, nil
}

func getNewNodeFunc(model *plugin.Plugin) node.NewNode{
	newFuncInterface, err := model.Lookup("NewNode")
	if err != nil{
		return nil
	}
	newFunc, ok := newFuncInterface.(node.NewNode)
	if !ok{
		return nil
	}
	return newFunc
}

func getNewTireManagerFunc(model *plugin.Plugin, path string) (trie.NewOptionTrie, error){
	newFuncInterface, err := model.Lookup("NewOptionTrie")
	if err != nil{
		return nil, errors.New("can not find NewOptionTrie from " + path)
	}
	newFunc, ok := newFuncInterface.(trie.NewOptionTrie)
	if !ok{
		return nil, errors.New("NewOptionTrie type is not trie.NewOptionTrie in " + path)
	}
	return newFunc, nil
}