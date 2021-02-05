package hash

import (
	"hash/crc32"
	"strconv"
	"strings"
)

func init(){
	lTime.Start()
}

func HashInit(num int) error{
	if num <= 0{
		return ErrDataBase
	}
	rootDataBase = make([]*Hash, num)
	t := lTime.GetTime()
	for i := 0; i < num; i++{
		n := &Hash{}
		n.SetKey(strconv.Itoa(i))
		n.SetTime(NodeTime{
			beginTime:t,
			durationTime:0,
		})
		rootDataBase[i] = n
	}
	return nil
}

func toHash (str string) int{
	hashNumber:=int(crc32.ChecksumIEEE([]byte(str)))
	if hashNumber>=0{
		return hashNumber%65536
	}
	return (-hashNumber)%65536
}

func toFatherList (str string) []string{
	return strings.Split(str,"/")
}

func toNode(num int,str string) (*Hash,error){
	if num >= len(rootDataBase) {
		return nil,ErrDataBase
	}
	if len(str) == 0{
		return rootDataBase[num],nil
	}
	strs := toFatherList(str)
	d := rootDataBase[num]
	for _,v := range strs{
		h,err := d.Get(v)
		if err != nil {
			return nil,err
		}
		if h.GetType() != "hash"{
			return nil,ErrNotHash
		}
		d=h.(*Hash)
	}
	return d,nil
}

func NewOtherNode(id,key,nodeType string,t int64) *Other{
	n := &Other{}
	n.SetKey(key)
	n.nodeType=nodeType
	n.id=id
	n.SetTime(NodeTime{
		beginTime:lTime.GetTime(),
		durationTime:t,
	})
	return n
}

func IsOverFlow(num int) error{
	if num < 0 || num >= len(rootDataBase){
		return ErrDataBase
	}
	return nil
}

func GetLength() int{
	return len(rootDataBase)
}

func Flush(num int) error{
	if num < 0 || num >= len(rootDataBase){
		return ErrDataBase
	}
	n := &Hash{}
	n.SetKey(strconv.Itoa(num))
	n.SetTime(NodeTime{
		beginTime:lTime.GetTime(),
		durationTime:0,
	})
	rootDataBase[num] = n
	return nil
}

func SetHash(num int,path,key string,t int64) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	n := &Hash{}
	n.SetKey(key)
	n.SetTime(NodeTime{
		beginTime:lTime.GetTime(),
		durationTime:t,
	})
	return h.Set(n)
}

func SetHashX(num int,path,key string,t int64) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	n := &Hash{}
	n.SetKey(key)
	n.SetTime(NodeTime{
		beginTime:lTime.GetTime(),
		durationTime:t,
	})
	h.SetX(n)
	return nil
}

func SetOther(num int,path,key,nodeType,id string,t int64) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	n := &Other{}
	n.SetKey(key)
	n.nodeType=nodeType
	n.id=id
	n.SetTime(NodeTime{
		beginTime:lTime.GetTime(),
		durationTime:t,
	})
	return h.Set(n)
}

func SetOtherX(num int,path,key,nodeType,id string,t int64) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	n := &Other{}
	n.SetKey(key)
	n.nodeType=nodeType
	n.id=id
	n.SetTime(NodeTime{
		beginTime:lTime.GetTime(),
		durationTime:t,
	})
	h.SetX(n)
	return nil
}

func GetID(num int,path,key string) (string,error){
	h,err := toNode(num,path)
	if err!=nil{
		return "",err
	}
	return h.GetNodeID(key)
}

func UpdateID(num int,path,key,id string) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	return h.UpdateNodeID(key,id)
}

func Select(num int,path string) ([]string,error){
	h,err := toNode(num,path)
	if err!=nil{
		return nil,err
	}
	return h.SelectKeyName(),nil
}

func SelectX(num int,path string) ([]Info,error){
	h,err := toNode(num,path)
	if err!=nil{
		return nil,err
	}
	return h.SelectInfo(),nil
}

func Delete(num int,path string,keys []string) (int,error){
	h,err := toNode(num,path)
	if err!=nil{
		return 0,err
	}
	return h.Del(keys),nil
}

func Type(num int,path,key string) (string,error){
	h,err := toNode(num,path)
	if err!=nil{
		return "",err
	}
	return h.GetNodeType(key)
}

func Exists(num int,path,key string) (bool,error){
	h,err := toNode(num,path)
	if err!=nil{
		return false,err
	}
	return h.Exists(key),nil
}

func Pex(num int,path,key string,t int64) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	return h.Pex(key,t)
}

func PexTo(num int,path,key string,t int64) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	return h.PexTo(key,t)
}

func Time(num int,path,key string) (int64,error){
	h,err := toNode(num,path)
	if err!=nil{
		return 0,err
	}
	return h.RTime(key)
}

func TimeTo(num int,path,key string) (int64,error){
	h,err := toNode(num,path)
	if err!=nil{
		return 0,err
	}
	return h.ETime(key)
}

func Rename(num int,path,key,newKey string) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	return h.Rename(key,newKey)
}

func RenameX(num int,path,key,newKey string) error{
	h,err := toNode(num,path)
	if err!=nil{
		return err
	}
	return h.RenameX(key,newKey)
}

func Size(num int,path string) (int,error){
	h,err := toNode(num,path)
	if err!=nil{
		return 0,err
	}
	return h.Size(),nil
}

func ToNode(num int,path string) (*Hash,error){
	return toNode(num,path)
}