package alertstate

import (
	//"encoding/json"
	"fmt"
)

type StateSlice struct {
	Total           StateUnit
	Genre           []StateTypeUnit `json:"type"`
	Sniffer         []StateSnifUnit `json:"sniffer"`
	SnifferRealtime []StateSnifUnit `json:"snifrealtime"`
	Site            []StateSiteUnit `json:"site"`
}

func (this *StateSlice) Init() {
	this.Genre = make([]StateTypeUnit, 0)
	this.Sniffer = make([]StateSnifUnit, 0)
	this.SnifferRealtime = make([]StateSnifUnit, 0)
	this.Site = make([]StateSiteUnit, 0)
}

type GlobalResult struct {
	Time int64
	StateSlice
}

func (this *GlobalResult) init() *GlobalResult {
	this.StateSlice.Init()
	return this
}

func (this *GlobalResult) String() string {
	str := fmt.Sprintln("time=", this.Time)
	str += fmt.Sprintln("Total=", this.Total)
	str += fmt.Sprintln("Genre=", this.Genre)
	str += fmt.Sprintln("Sniffer=", this.Sniffer)
	str += fmt.Sprintln("site=", this.Site)
	return str
}

type StateSiteUnits []StateSiteUnit

func (this StateSiteUnits) GenOneSite(id int32, unit *StateUnit) *StateSiteUnit {
	return &StateSiteUnit{
		Id: uint32(id),
		StateTypeUnit: StateTypeUnit{
			Name:      gIdNameMap.GetSiteName(id),
			StateUnit: *unit,
		},
	}
}

func (this StateSiteUnits) FromMap(siteNum map[int32]*StateUnit) []StateSiteUnit {
	ss := []StateSiteUnit(this)
	for siteid, state := range siteNum {
		s := this.GenOneSite(siteid, state)
		ss = append(ss, *s)
	}
	return ss
}

type StateSiteUnit struct {
	Id uint32 `json:"id"`
	StateTypeUnit
}

func (this *StateSiteUnit) String() string {
	str := fmt.Sprintf("id:%d", this.Id)
	str += fmt.Sprintln(this.StateTypeUnit)
	return str
}

type StateSnifUnits []StateSnifUnit

func (this StateSnifUnits) GenOneType(id int32, unit *StateUnit) *StateTypeUnit {
	typestr, ok := classIntStr[classType(id)]
	if !ok {
		panic(fmt.Sprintln("StateSnifUnits.GenOneType classIntStr id not exist ", id))
	}
	return &StateTypeUnit{
		Name:      typestr,
		StateUnit: *unit,
	}
}

func (this StateSnifUnits) GenOneSniffer(id int32, unit *StateUnit, types []StateTypeUnit) *StateSnifUnit {
	return &StateSnifUnit{
		Id: uint32(id),
		//Id: gIdNameMap.GetSnifferName(id),
		StateTypeUnit: StateTypeUnit{
			Name: gIdNameMap.GetSnifferName(id),
			StateUnit: StateUnit{
				Total:  unit.Total,
				Noread: unit.Noread,
			},
		},
		Types: types,
	}
}

func (this StateSnifUnits) FromMap(snifferNum map[int32]*StateUnit,
	snifTypeNum map[int32]map[int32]*StateUnit) []StateSnifUnit {

	sss := []StateSnifUnit(this)

	for snifid, SnifStateUnit := range snifferNum {
		ts := make([]StateTypeUnit, 0, 5)
		if typenum, ok := snifTypeNum[snifid]; ok {
			for typeid, num := range typenum {
				t := this.GenOneType(typeid, num)
				ts = append(ts, *t)
			}
		} else {
			//panic(fmt.Sprintln("StateSnifUnits.FromMap sniffer id not int snifTypeNum ", snifid))
		}
		sniffer := this.GenOneSniffer(snifid, SnifStateUnit, ts)
		sss = append(sss, *sniffer)
	}
	return sss
}

type StateSnifUnit struct {
	Id uint32 `json:"id"`
	StateTypeUnit
	Types []StateTypeUnit `json:"types"`
}

func (this *StateSnifUnit) String() string {
	str := fmt.Sprintf("id=%d", this.Id)
	str += fmt.Sprintln(this.StateTypeUnit)
	str += fmt.Sprintln(this.Types)
	return str
}

/*
func (this []StateTypeUnit) FromMap(typeNum map[int32]*StateUnit) {
	//ss := []StateTypeUnit(this)
	for id, stateUnit := range typeNum {
		s := StateTypeUnits(this).GenOneType(id, stateUnit)
		this = append(this, *s)
	}
	fmt.Println("ss=", this)
	this = StateTypeUnits(this)
	fmt.Println("this=", this)
}
*/
type StateTypeUnits []StateTypeUnit

func (this StateTypeUnits) GenOneType(id int32, state *StateUnit) *StateTypeUnit {
	typestr, ok := classIntStr[classType(id)]
	if !ok {
		panic(fmt.Sprintln("StateTypeUnits.FromMap id not exist ", id))
	}
	return &StateTypeUnit{
		Name:      typestr,
		StateUnit: *state,
	}
}

func (this StateTypeUnits) FromMap(typeNum map[int32]*StateUnit) []StateTypeUnit {
	ss := []StateTypeUnit(this)
	for id, stateUnit := range typeNum {
		s := this.GenOneType(id, stateUnit)
		ss = append(ss, *s)
	}
	//fmt.Println("ss=", ss)
	return ss
}

type StateTypeUnit struct {
	//id   uint32
	Name string `json:"name"`
	StateUnit
}

func (this StateTypeUnit) String() string {
	str := fmt.Sprintf("name:%s", this.Name)
	str += fmt.Sprintln(&this.StateUnit)
	return str
}

type StateUnit struct {
	Total  int32 `json:"total"`
	Noread int32 `json:"noread"`
}

func (this *StateUnit) String() string {
	return fmt.Sprintf("total:%d;noread:%d", this.Total, this.Noread)
}

func (this *StateUnit) merge(para *StateUnit) {
	this.Total += para.Total
	this.Noread += para.Noread
}
