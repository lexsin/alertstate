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
	snifTypeNumRt map[int64]map[int32]*StateUnit //realtime
	snifferNumRt  map[int64]*StateUnit
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

func (this *GlobalCache) Insert(data EntryRecord) {
	this.lock.Lock()
	defer this.lock.Unlock()
	//to global_cache
	this.StateMap.accmAdd(data)
	//record name
	gIdNameMap.Insert(
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
	result.Sniffer = StateSnifUnits(result.Sniffer).FromMap(this.snifferNum, this.snifTypeNum)
	result.Site = StateSiteUnits(result.Site).FromMap(this.siteNum)
	result.SnifferRealtime = StateSnifUnits(result.SnifferRealtime).FromMap(this.snifferNumRt, this.snifTypeNumRt)
	this.snifferNumRt = nil
	this.snifTypeNumRt = nil
	return result
}

func SubType(genre string, total int32, noread int32) {
	id := int32(classStrInt[genre])
	gGlobalCache.lock.Lock()
	defer gGlobalCache.lock.Unlock()
	state, ok := gGlobalCache.typeNum[id]
	if ok {
		state.Total -= total
		if state.Total <= 0 {
			state.Total = 0
		}
		state.Noread -= noread
		if state.Noread <= 0 {
			state.Noread = 0
		}
	}
}

func AddType(genre string, total int32, noread int32) {
	id := int32(classStrInt[genre])
	gGlobalCache.lock.Lock()
	defer gGlobalCache.lock.Unlock()
	state, ok := gGlobalCache.typeNum[id]
	if ok {
		state.Total += total
		state.Noread += noread
	} else if !ok {
		gGlobalCache.typeNum[id] = &StateUnit{
			Total:  total,
			Noread: noread,
		}
	}
}

func AddTotal(total int32, noread int32) {
	gGlobalCache.lock.Lock()
	defer gGlobalCache.lock.Unlock()
	gGlobalCache.total.Total += total
	gGlobalCache.total.Noread += noread
}

func SubTotal(total int32, noread int32) {
	gGlobalCache.lock.Lock()
	defer gGlobalCache.lock.Unlock()
	gGlobalCache.total.Total -= total
	if gGlobalCache.total.Total < 0 {
		gGlobalCache.total.Total = 0
	}
	gGlobalCache.total.Noread -= noread
	if gGlobalCache.total.Noread < 0 {
		gGlobalCache.total.Noread = 0
	}
}

func AddSniffer(id int64, name string, total int32, noread int32) {
	//inset name
	gIdNameMap.InsertSniffer(idNameT{int64(id), name})
	//insert id
	gGlobalCache.lock.Lock()
	defer gGlobalCache.lock.Unlock()
	state, ok := gGlobalCache.snifferNum[id]
	if ok {
		state.Total += total
		state.Noread += noread
	} else if !ok {
		gGlobalCache.snifferNum[id] = &StateUnit{
			Total:  total,
			Noread: noread,
		}
	}
}

func SubSniffer(id int64, total int32, noread int32) {
	gGlobalCache.lock.Lock()
	defer gGlobalCache.lock.Unlock()
	state, ok := gGlobalCache.snifferNum[id]
	if ok {
		state.Total -= total
		if state.Total == 0 {
			state.Total = 0
		}
		state.Noread -= noread
		if state.Noread == 0 {
			state.Noread = 0
		}
	}
}

func AddSite(id int32, name string, total int32, noread int32) {
	//insert name
	gIdNameMap.InsertSite(idNameT{int64(id), name})
	//insert id
	gGlobalCache.lock.Lock()
	defer gGlobalCache.lock.Unlock()
	state, ok := gGlobalCache.siteNum[int64(id)]
	if ok {
		state.Total += total
		state.Noread += noread
	} else if !ok {
		gGlobalCache.siteNum[int64(id)] = &StateUnit{
			Total:  total,
			Noread: noread,
		}
	}
}

func SubSite(id int32, total int32, noread int32) {
	gGlobalCache.lock.Lock()
	defer gGlobalCache.lock.Unlock()
	state, ok := gGlobalCache.siteNum[int64(id)]
	if ok {
		state.Total -= total
		if state.Total == 0 {
			state.Total = 0
		}
		state.Noread -= noread
		if state.Noread == 0 {
			state.Noread = 0
		}
	}
}
