package alertstate

import (
	"fmt"
	"sync"
)

var gIdNameMap IdNameMap

type idNameT struct {
	id   int64
	name string
}

type IdNameMap struct {
	lock    *sync.Mutex
	sniffer map[int64]string
	site    map[int64]string
}

func (this *IdNameMap) Init() *IdNameMap {
	this.lock = new(sync.Mutex)
	this.sniffer = make(map[int64]string)
	this.site = make(map[int64]string)
	return this
}

func (this *IdNameMap) InsertSniffer(sniffer idNameT) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if sniffer.id != 0 {
		this.sniffer[sniffer.id] = sniffer.name
	}
}

func (this *IdNameMap) InsertSite(site idNameT) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.site[site.id] = site.name
}

func (this *IdNameMap) Insert(sniffer idNameT, site idNameT) {
	this.InsertSniffer(sniffer)
	this.InsertSite(site)
}

func (this *IdNameMap) GetSnifferName(id int64) string {
	name, ok := this.sniffer[id]
	if !ok {
		panic(fmt.Sprintln("sniffer name not exist with id", id))
	}
	return name
}

func (this *IdNameMap) GetSiteName(id int64) string {
	name, ok := this.site[id]
	if !ok {
		panic(fmt.Sprintln("site name not exist with id", id))
	}
	return name
}
