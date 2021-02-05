package localTime

import "time"

func InitTime() *Timer{
	return &Timer{
		timer:time.Now().Unix(),
	}
}

func (handle *Timer) Start(){
	go handle.run()
	go handle.correctionTimer()
}

func (handle *Timer) run(){
	for {
		time.Sleep(time.Second)
		handle.AddTime()
	}
}

func (handle *Timer) correctionTimer() {
	for {
		time.Sleep(time.Hour)
		handle.SetTime(time.Now().Unix())
	}
}

func (handle *Timer) GetTime() int64{
	handle.lock.RLock()
	defer handle.lock.RUnlock()
	return handle.timer
}

func (handle *Timer) SetTime(t int64) {
	handle.lock.Lock()
	defer handle.lock.Unlock()
	handle.timer=t
}

func (handle *Timer) AddTime() {
	handle.lock.Lock()
	defer handle.lock.Unlock()
	handle.timer++
}