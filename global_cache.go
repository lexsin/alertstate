package alertstate

import (
	"fmt"
	"sync"
)

//var gGlobalCache GlobalCache

type GlobalCache struct {
	lock *sync.Mutex
	time int64
	StateMap
	snifTypeNumRt map[int64]map[int32]*StateUnit //realtime
	snifferNumRt  map[int64]*StateUnit
	IdNameMap
}

func (this *GlobalCache) String() string {
	str := fmt.Sprintln("time:", this.time)
	str += fmt.Sprintln("statemap:", &this.StateMap)
	str += fmt.Sprintln("snifTypeNumRt:", &this.snifTypeNumRt)
	return str
}

func (this *GlobalCache) Init() {
	this.lock = new(sync.Mutex)
	this.StateMap.Init()
	this.IdNameMap.Init()
}

func (this *GlobalCache) Insert(data EntryRecord) {
	this.lock.Lock()
	defer this.lock.Unlock()
	//to global_cache
	this.StateMap.accmAdd(data)
	//record name
	this.IdNameMap.Insert(
		idNameT{int64(data.Sniffer), data.Sniffername},
		idNameT{int64(data.Site), data.Sitename})
}

func (this *GlobalCache) ToSlice() *GlobalResult {
	this.lock.Lock()
	defer this.lock.Unlock()

	result := new(GlobalResult).init()
	result.Time = this.time
	result.Total = this.total
	result.Genre = StateTypeUnits(result.Genre).FromMap(this.typeNum)
	//fmt.Println("genre=", result.Genre)
	result.Sniffer = StateSnifUnits(result.Sniffer).FromMap(this.snifferNum, this.snifTypeNum, &this.IdNameMap)
	result.Site = StateSiteUnits(result.Site).FromMap(this.siteNum, &this.IdNameMap)
	result.SnifferRealtime = StateSnifUnits(result.SnifferRealtime).FromMap(this.snifferNumRt, this.snifTypeNumRt, &this.IdNameMap)
	this.snifferNumRt = nil
	this.snifTypeNumRt = nil
	return result
}
