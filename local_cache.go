package alertstate

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

type EntryRecord struct {
	Timestamp   int64
	Genre       string
	genre       classType
	Sniffer     int64
	Sniffername string
	Site        int32
	Sitename    string
}

func (this *EntryRecord) Transform() error {
	genre, ok := classStrInt[this.Genre]
	if !ok {
		//fmt.Println("alert type error:", this.Genre)
		Error("alert type error:", this.Genre)
		return ErrClassTypeErr
	}
	this.genre = genre
	return nil
}

//var GLocalCach LocalCache

//var InputCh = make(chan EntryRecord, InputCacheLenDef)

type LocalCache struct {
	//AlertState
	lock        *sync.Mutex
	globalCache GlobalCache
	WinWidth    int32 //秒
	WinNum      int32
	InputCh     chan EntryRecord
	MaxWinId    int64
	Windows     map[int64]*window
	WillBeDelId int64
	Handler     func()
	MvLocToGlb  func(id int64) error
	Done        chan int8
}

func (this *LocalCache) Init(width int32, num int32) (error, *LocalCache) {
	this.lock = new(sync.Mutex)
	if width <= 0 || num < 1 {
		return ErrWinParamentErr, nil
	}
	this.WinWidth = width
	this.WinNum = num
	this.InputCh = make(chan EntryRecord, InputCacheLenDef)
	this.Windows = make(map[int64]*window)
	this.Done = make(chan int8, 1)
	return nil, this
}

func (this *LocalCache) GetOvertimeWinid() int64 {
	return int64(int64(this.MaxWinId) - int64(this.WinNum))
}

var deleteWindowCount = 0

func (this *LocalCache) deleteWindow(id int64) error {
	Debug("delete window id=", id)
	//move to global cache
	this.MvToGlobal(id)
	//to handler(websocket)
	//TODO
	this.Handler()
	return nil
}

func (this *LocalCache) recycleAllWindows() {
	for wid, _ := range this.Windows {
		fmt.Println("delete window id=", wid)
		this.MvToGlobal(wid)
	}
	this.Handler()
	this.MaxWinId = 0
}

// func (this *LocalCache) StateNHandler() {
// 	result := gGlobalCache.ToSlice()
// 	this.Handler(result)
// 	//fmt.Println("result", result)
// }

func (this *LocalCache) MvToGlobal(id int64) error {
	this.globalCache.lock.Lock()
	defer this.globalCache.lock.Unlock()

	if _, ok := this.Windows[id]; !ok {
		deleteWindowCount++
		Info("delete window not exist count=", deleteWindowCount)
		return ErrWindowNotExist
	}

	this.globalCache.time = this.Windows[id].time
	this.globalCache.merge((this.Windows[id].mp))
	this.globalCache.snifTypeNumRt = this.Windows[id].mp.snifTypeNum
	this.globalCache.snifferNumRt = this.Windows[id].mp.snifferNum
	//free map
	delete(this.Windows, id)
	return nil
}

func (this *LocalCache) insert(winid int64, data EntryRecord) {
	//record name
	this.globalCache.IdNameMap.Insert(
		idNameT{int64(data.Sniffer), data.Sniffername},
		idNameT{int64(data.Site), data.Sitename})
	//insert
	this.Windows[winid].insert(data)
}

func deferf(winid int64, maxid int64, time int64) {
	if err := recover(); err != nil {
		Error("panic:", err, string(debug.Stack()))
		Error("winid=", winid, "maxid=", maxid, "time=", time)
	}
}

func (this *LocalCache) Insert(data EntryRecord) (err error) {
	winid := data.Timestamp / int64(this.WinWidth)
	defer deferf(winid, this.MaxWinId, data.Timestamp)
	if this.MaxWinId == 0 {
		/*
			new window
			value maxid
			insert
		*/
		this.Windows[winid] = new(window).new(winid, this.WinWidth)
		this.MaxWinId = winid
	} else if winid == this.MaxWinId-1 {
		/*
			new window if not exist
			insert
		*/
		if _, ok := this.Windows[winid]; !ok {
			this.Windows[winid] = new(window).new(winid, this.WinWidth)
		}
	} else if winid == this.MaxWinId {
		/*
			insert
		*/
	} else if winid == this.MaxWinId+1 {
		/*
			new window
			value maxid
			go delete old window
			insert
		*/
		this.Windows[winid] = new(window).new(winid, this.WinWidth)
		this.deleteWindow(this.MaxWinId - 1)
		this.MaxWinId = winid
	} else if winid < this.MaxWinId-1 {
		/*
			insert global window
		*/
		this.globalCache.Insert(data)
		return
	} else if winid > this.MaxWinId+1 {
		/*
			delete all window
			new window
			value maxid
			insert
		*/
		for wid, _ := range this.Windows {
			this.MvToGlobal(wid)
		}
		this.Windows[winid] = new(window).new(winid, this.WinWidth)
		this.MaxWinId = winid
	}
	this.insert(winid, data)
	return nil
}

func (this *LocalCache) Start() {
	if this.Handler == nil {
		fmt.Println("start Error:LocalCache handler is nil")
		return
	}

	winTimeout := time.Duration(this.WinWidth) * time.Second
	timer := time.NewTimer(winTimeout)
	for {
		timer.Reset(winTimeout)
		select {
		case data := <-this.InputCh:
			if err := data.Transform(); err != nil {
				continue
			}
			this.Insert(data)
		case <-timer.C:
			this.recycleAllWindows()
		case <-this.Done:
			return
		}
	}
}

func Print() {
	//	go func() {
	//		for {
	//			time.Sleep(1 * time.Second)

	//			fmt.Println("snifTypeNumRt=", gGlobalCache.snifTypeNumRt, "snifferNumRt=", gGlobalCache.snifTypeNumRt)
	//		}
	//	}()
}
