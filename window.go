package alertstate

import (
	"fmt"
)

type StateMap struct {
	typeNum     map[int32]*StateUnit
	snifferNum  map[int32]*StateUnit
	snifTypeNum map[int32]map[int32]*StateUnit
	siteNum     map[int32]*StateUnit
	total       StateUnit
}

func (this *StateMap) String() string {
	str := fmt.Sprintf("mp{")
	str += fmt.Sprintln("typeNum:", this.typeNum)
	str += fmt.Sprintln("snifferNum:", this.snifferNum)
	str += fmt.Sprintln("snifTypeNum:", this.snifTypeNum)
	str += fmt.Sprintln("siteNum:", this.siteNum)
	str += fmt.Sprintln("total:", this.total)
	str += fmt.Sprintf("}")
	return str
}

func (this *StateMap) accmAdd(data EntryRecord) {
	MapInt32StateUnit(this.typeNum).Add(int32(data.genre))
	MapInt32StateUnit(this.snifferNum).Add(int32(data.Sniffer))
	MapInt32Int32StateUnit(this.snifTypeNum).Add(int32(data.Sniffer), int32(data.genre))
	MapInt32StateUnit(this.siteNum).Add(data.Site)
	this.total.Total++
	this.total.Noread++
}

func (this *StateMap) merge(para StateMap) error {
	this.total.merge(&(para.total))
	MapInt32StateUnit(this.typeNum).merge(MapInt32StateUnit(para.typeNum))
	MapInt32StateUnit(this.snifferNum).merge(MapInt32StateUnit(para.snifferNum))
	MapInt32StateUnit(this.siteNum).merge(MapInt32StateUnit(para.siteNum))
	MapInt32Int32StateUnit(this.snifTypeNum).merge(MapInt32Int32StateUnit(para.snifTypeNum))
	return nil
}

func (this *StateMap) Init() *StateMap {
	this.typeNum = make(map[int32]*StateUnit)
	this.snifferNum = make(map[int32]*StateUnit)
	this.snifTypeNum = make(map[int32]map[int32]*StateUnit)
	this.siteNum = make(map[int32]*StateUnit)
	return this
}

type window struct {
	id   int32
	time int64
	//state State
	mp StateMap
}

func (this *window) String() string {
	str := fmt.Sprintf("id=%d time=%d", this.id, this.time)
	str += fmt.Sprintln(&this.mp)
	return str
}

func (this *window) init(timestamp int64) *window {
	this.time = timestamp
	this.mp = *(new(StateMap).Init())
	return this
}

func (this *window) insert(data EntryRecord) {
	this.mp.accmAdd(data)
}

type MapInt32Int32StateUnit map[int32]map[int32]*StateUnit

func (this MapInt32Int32StateUnit) Add(key1 int32, key2 int32) {
	mmp := map[int32]map[int32]*StateUnit(this)
	_, ok := mmp[key1]
	if !ok {
		mmp[key1] = make(map[int32]*StateUnit)
	}
	MapInt32StateUnit(mmp[key1]).Add(key2)
}

func (this MapInt32Int32StateUnit) merge(para MapInt32Int32StateUnit) {
	//	src := map[int32]map[int32]*StateUnit(this)
	//	dst := map[int32]map[int32]*StateUnit(para)
	//	for srck, srcv := range src {
	//	}
}

type MapInt32StateUnit map[int32]*StateUnit

func (this MapInt32StateUnit) Add(key int32) {
	mp := map[int32]*StateUnit(this)
	p, ok := mp[key]
	if ok {
		p.Total += 1
		p.Noread += 1
	} else {
		mp[key] = &StateUnit{
			Total:  1,
			Noread: 1,
		}
	}
}

func (this MapInt32StateUnit) merge(para MapInt32StateUnit) {
	dst := map[int32]*StateUnit(this)
	src := map[int32]*StateUnit(para)
	for srck, srcv := range src {
		if dstv, ok := dst[srck]; ok {
			dstv.merge(srcv)
		} else {
			dst[srck] = &StateUnit{
				Total:  srcv.Total,
				Noread: srcv.Noread,
			}
		}
	}
}
