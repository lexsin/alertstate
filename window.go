package alertstate

import (
	"fmt"
)

type StateMap struct {
	typeNum     map[int32]*StateUnit
	snifferNum  map[int64]*StateUnit
	snifTypeNum map[int64]map[int32]*StateUnit
	siteNum     map[int64]*StateUnit
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
	MapInt64StateUnit(this.snifferNum).Add(int64(data.Sniffer))
	MapInt64Int32StateUnit(this.snifTypeNum).Add(int64(data.Sniffer), int32(data.genre))
	MapInt64StateUnit(this.siteNum).Add(data.Site)
	this.total.Total++
	this.total.Noread++
}

func (this *StateMap) merge(para StateMap) error {
	this.total.merge(&(para.total))
	MapInt32StateUnit(this.typeNum).merge(MapInt32StateUnit(para.typeNum))
	MapInt64StateUnit(this.snifferNum).merge(MapInt64StateUnit(para.snifferNum))
	MapInt64StateUnit(this.siteNum).merge(MapInt64StateUnit(para.siteNum))
	MapInt64Int32StateUnit(this.snifTypeNum).merge(MapInt64Int32StateUnit(para.snifTypeNum))
	return nil
}

func (this *StateMap) Init() *StateMap {
	this.typeNum = make(map[int32]*StateUnit)
	this.snifferNum = make(map[int64]*StateUnit)
	this.snifTypeNum = make(map[int64]map[int32]*StateUnit)
	this.siteNum = make(map[int64]*StateUnit)
	return this
}

type window struct {
	id   int64
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

type MapInt64Int32StateUnit map[int64]map[int32]*StateUnit

func (this MapInt64Int32StateUnit) Add(key1 int64, key2 int32) {
	mmp := map[int64]map[int32]*StateUnit(this)
	_, ok := mmp[key1]
	if !ok {
		mmp[key1] = make(map[int32]*StateUnit)
	}
	MapInt32StateUnit(mmp[key1]).Add(key2)
}

func (this MapInt64Int32StateUnit) merge(para MapInt64Int32StateUnit) {
	//	src := map[int64]map[int64]*StateUnit(this)
	//	dst := map[int64]map[int64]*StateUnit(para)
	//	for srck, srcv := range src {
	//	}
}

type MapInt64StateUnit map[int64]*StateUnit

func (this MapInt64StateUnit) Add(key int64) {
	mp := map[int64]*StateUnit(this)
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
func (this MapInt64StateUnit) merge(para MapInt64StateUnit) {
	dst := map[int64]*StateUnit(this)
	src := map[int64]*StateUnit(para)
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
