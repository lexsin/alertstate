package alertstate

type AlertState struct {
	//globalCache GlobalCache
	//gIdNameMap IdNameMap
	LocalCach LocalCache
	InputCh   chan EntryRecord
	Handler   func(*GlobalResult)
}

func (this *AlertState) Init(width int32) *AlertState {
	//initLog()
	//WindowWidth := 5
	//this.globalCache.Init()
	//this.gIdNameMap.Init()
	this.LocalCach.Init(int32(width), 2)
	this.LocalCach.globalCache.Init()
	this.InputCh = this.LocalCach.InputCh
	//this.nameMapP.Init()
	return this
}

func (this *AlertState) localCacheHandler() {
	result := this.LocalCach.globalCache.ToSlice()
	this.Handler(result)
}

func (this *AlertState) Start(handler func(*GlobalResult)) error {
	this.Handler = handler
	this.LocalCach.Handler = this.localCacheHandler
	this.LocalCach.Start()
	return nil
}
func (this *AlertState) Free() {
	this.LocalCach.Done <- 1
}

func (this *AlertState) WinWidth() int32 {
	return this.LocalCach.WinWidth
}

// func (this *AlertState) MvLocToGlb(id int64) error {
// 	this.gLocalCach.lock.Lock()
// 	defer this.gLocalCach.lock.Unlock()

// 	if _, ok := this.gLocalCach.Windows[id]; !ok {
// 		deleteWindowCount++
// 		Info("delete window not exist count=", deleteWindowCount)
// 		return ErrWindowNotExist
// 	}

// 	this.gLocalCach.globalCache.time = this.gLocalCach.Windows[id].time
// 	this.gLocalCach.globalCache.merge((this.gLocalCach.Windows[id].mp))
// 	this.gLocalCach.snifTypeNumRt = this.gLocalCach.Windows[id].mp.snifTypeNum
// 	this.gLocalCach.snifferNumRt = this.gLocalCach.Windows[id].mp.snifferNum
// 	//free map
// 	delete(this.gLocalCach.Windows, id)
// 	return nil
// }

func (this *AlertState) SubType(genre string, total int32, noread int32) {
	id := int32(classStrInt[genre])
	this.LocalCach.globalCache.lock.Lock()
	defer this.LocalCach.globalCache.lock.Unlock()
	state, ok := this.LocalCach.globalCache.typeNum[id]
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

func (this *AlertState) AddType(genre string, total int32, noread int32) {
	id := int32(classStrInt[genre])
	this.LocalCach.globalCache.lock.Lock()
	defer this.LocalCach.globalCache.lock.Unlock()
	state, ok := this.LocalCach.globalCache.typeNum[id]
	if ok {
		state.Total += total
		state.Noread += noread
	} else if !ok {
		this.LocalCach.globalCache.typeNum[id] = &StateUnit{
			Total:  total,
			Noread: noread,
		}
	}
}

func (this *AlertState) AddTotal(total int32, noread int32) {
	this.LocalCach.globalCache.lock.Lock()
	defer this.LocalCach.globalCache.lock.Unlock()
	this.LocalCach.globalCache.total.Total += total
	this.LocalCach.globalCache.total.Noread += noread
}

func (this *AlertState) SubTotal(total int32, noread int32) {
	this.LocalCach.globalCache.lock.Lock()
	defer this.LocalCach.globalCache.lock.Unlock()
	this.LocalCach.globalCache.total.Total -= total
	if this.LocalCach.globalCache.total.Total < 0 {
		this.LocalCach.globalCache.total.Total = 0
	}
	this.LocalCach.globalCache.total.Noread -= noread
	if this.LocalCach.globalCache.total.Noread < 0 {
		this.LocalCach.globalCache.total.Noread = 0
	}
}

func (this *AlertState) AddSniffer(id int64, name string, total int32, noread int32) {
	//inset name
	this.LocalCach.globalCache.InsertSniffer(idNameT{int64(id), name})
	//insert id
	this.LocalCach.globalCache.lock.Lock()
	defer this.LocalCach.globalCache.lock.Unlock()
	state, ok := this.LocalCach.globalCache.snifferNum[id]
	if ok {
		state.Total += total
		state.Noread += noread
	} else if !ok {
		this.LocalCach.globalCache.snifferNum[id] = &StateUnit{
			Total:  total,
			Noread: noread,
		}
	}
}

func (this *AlertState) SubSniffer(id int64, total int32, noread int32) {
	this.LocalCach.globalCache.lock.Lock()
	defer this.LocalCach.globalCache.lock.Unlock()
	state, ok := this.LocalCach.globalCache.snifferNum[id]
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

func (this *AlertState) AddSite(id int32, name string, total int32, noread int32) {
	//insert name
	this.LocalCach.globalCache.InsertSite(idNameT{int64(id), name})
	//insert id
	this.LocalCach.globalCache.lock.Lock()
	defer this.LocalCach.globalCache.lock.Unlock()
	state, ok := this.LocalCach.globalCache.siteNum[int64(id)]
	if ok {
		state.Total += total
		state.Noread += noread
	} else if !ok {
		this.LocalCach.globalCache.siteNum[int64(id)] = &StateUnit{
			Total:  total,
			Noread: noread,
		}
	}
}

func (this *AlertState) SubSite(id int32, total int32, noread int32) {
	this.LocalCach.globalCache.lock.Lock()
	defer this.LocalCach.globalCache.lock.Unlock()
	state, ok := this.LocalCach.globalCache.siteNum[int64(id)]
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
