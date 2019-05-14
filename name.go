package alertstate

import (
	"fmt"
)

var gIdNameMap IdNameMap

type idNameT struct {
	id   int32
	name string
}

type IdNameMap struct {
	sniffer map[int32]string
	site    map[int32]string
}

func (this *IdNameMap) Init() *IdNameMap {
	this.sniffer = make(map[int32]string)
	this.site = make(map[int32]string)
	return this
}

func (this *IdNameMap) Insert(sniffer idNameT, site idNameT) {
	if sniffer.id != 0 {
		this.sniffer[sniffer.id] = sniffer.name
	}
	if site.id != 0 {
		this.site[site.id] = site.name
	}
}

func (this *IdNameMap) GetSnifferName(id int32) string {
	name, ok := this.sniffer[id]
	if !ok {
		panic(fmt.Sprintln("sniffer name not exist with id", id))
	}
	return name
}

func (this *IdNameMap) GetSiteName(id int32) string {
	name, ok := this.site[id]
	if !ok {
		panic(fmt.Sprintln("site name not exist with id", id))
	}
	return name
}
