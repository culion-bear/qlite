package persistence

import "time"

const(
	ERROR	=	" (ERROR) "
	INFO	=	" (INFO) "
	WARNING	=	" (WARNING) "
)

const(
	SEND	=	" -SEND- "
	READ	=	" -READ- "
	SAVE	=	" -SAVE- "
)

type LogManager struct {
	writer *WriteManager
}

func NewLogHandle(fileName string,p func(string,error)) *LogManager{
	v := &LogManager{
		writer:NewWriter(fileName,0),
	}
	v.writer.Start(p)
	return v
}

func (handle *LogManager) Close(){
	handle.writer.Close()
}

func (handle *LogManager) Write(msg,level,action string){
	handle.writer.Write(time.Now().Format("[2006-01-02 15:04:05]")+action+level+msg+"\n")
}