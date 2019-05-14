package alertstate

import (
	"fmt"
	"sync"
)

var gGlobalCache GlobalCache

type GlobalCache struct {
	lock *sync.Mutex
	time int64
	StateMap
	snifTypeNumRt map[int32]map[int32]*StateUnit //realtime
	snifferNumRt  map[int32]*StateUnit
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
}

func (this *GlobalCache) Insert(data entryRecord) {
	this.lock.Lock()
	defer this.lock.Unlock()
	//to global_cache
	this.StateMap.accmAdd(data)
	//record name
	gIdNameMap.Insert(
		idNameT{data.sniffer, data.sniffername},
		idNameT{data.site, data.sitename})
}

func (this *GlobalCache) ToSlice() *GlobalResult {
	this.lock.Lock()
	defer this.lock.Unlock()

	result := new(GlobalResult).init()
	result.Time = this.time
	result.Total = this.total
	result.Class = StateTypeUnits(result.Class).FromMap(this.typeNum)
	//fmt.Println("class=", result.Class)
	result.Sniffer = StateSnifUnits(result.Sniffer).FromMap(this.snifferNum, this.snifTypeNum)
	result.Site = StateSiteUnits(result.Site).FromMap(this.siteNum)
	result.SnifferRealtime = StateSnifUnits(result.Sniffer).FromMap(this.snifferNumRt, this.snifTypeNumRt)
	return result
}
