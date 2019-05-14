package main

import (
	"fmt"
	"ifmes/modules/alertstate"
)

func alert_state_handler(result *alertstate.GlobalResult) {
	//fmt.Println(result)
	fix := AlertStateFix{}
	fix.FromArts(result)
	fmt.Println("fix=", fix)
	//to websockt

	sites := AlertStateSite{}
	sites.FromArts(result)
	fmt.Println("sites=", sites)
	//to websockt

	realtimesniffer := AlertRealtimeSniffer{}
	realtimesniffer.FromArts(result)
	fmt.Println("realtimesniffer=", realtimesniffer)
	//to websockt
}

type StateSnifUnit struct {
	Id int32
	alertstate.StateTypeUnit
}

func (this *StateSnifUnit) FromArts(src *alertstate.StateSnifUnit) {
	this.Id = int32(src.Id)
	this.Name = src.Name
	this.Total = src.Total
	this.Noread = src.Noread
}

type AlertStateFix struct {
	Time     string
	Interval int32
	Total    alertstate.StateUnit       `json:"total"`
	Genre    []alertstate.StateTypeUnit `json:"type"`
	Sniffer  []StateSnifUnit            `json:"sniffer"`
}

func (this *AlertStateFix) FromArts(para *alertstate.GlobalResult) {
	this.Time = TimeFormUi(para.Time)
	this.Interval = alertstate.GLocalCach.WinWidth * 1000 //s-->ms
	this.Total = para.Total
	this.Genre = para.Class
	this.SnifferFromArts(para.Sniffer)
}

func (this *AlertStateFix) SnifferFromArts(para []alertstate.StateSnifUnit) {
	this.Sniffer = make([]StateSnifUnit, 0)
	for _, snif := range para {
		sniff := StateSnifUnit{}
		sniff.FromArts(&snif)
		this.Sniffer = append(this.Sniffer, sniff)
	}
}

type RealTimeSnifferUnit struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	//Count int32
	Realtime RealTimeSnifferUnitSub `json:"realtime"`
}

func (this *RealTimeSnifferUnit) FromArts(src *alertstate.StateSnifUnit) {
	this.Id = int32(src.Id)
	this.Name = src.Name
	this.Realtime.FromArts(src)
}

type RealTimeSnifferUnitSub struct {
	Count int32              `json:"count"`
	Types []RealTimeTypeUnit `json:"types"`
}

func (this *RealTimeSnifferUnitSub) FromArts(src *alertstate.StateSnifUnit) {
	this.Count = src.Total
	this.Types = make([]RealTimeTypeUnit, 0)
	for _, genre := range src.Types {
		t := RealTimeTypeUnit{}
		t.FromArts(&genre)
		this.Types = append(this.Types, t)
	}
}

type RealTimeTypeUnit struct {
	Name  string `json:"name"`
	Count int32  `json:"count"`
}

func (this *RealTimeTypeUnit) FromArts(src *alertstate.StateTypeUnit) {
	this.Name = src.Name
	this.Count = src.Total
}

type AlertRealtimeSniffer struct {
	Time     string                `json:"time"`
	Interval int32                 `json:"interval"`
	Sniffer  []RealTimeSnifferUnit `json:"sniffer"`
}

func (this *AlertRealtimeSniffer) SnifferFromArts(src []alertstate.StateSnifUnit) {
	this.Sniffer = make([]RealTimeSnifferUnit, 0)
	for _, srcSniffer := range src {
		sniff := RealTimeSnifferUnit{}
		sniff.FromArts(&srcSniffer)
		this.Sniffer = append(this.Sniffer, sniff)
	}
}

func (this *AlertRealtimeSniffer) FromArts(src *alertstate.GlobalResult) {
	this.Time = TimeFormUi(src.Time)
	this.Interval = alertstate.GLocalCach.WinWidth * 1000
	this.SnifferFromArts(src.SnifferRealtime)
}

type AlertStateSiteunit struct {
	Id     int32
	Name   string
	Total  int32
	Noread int32
}

func (this *AlertStateSiteunit) FromArts(src *alertstate.StateSiteUnit) {
	this.Id = int32(src.Id)
	this.Name = src.Name
	this.Total = src.Total
	this.Noread = src.Noread
}

type AlertStateSite struct {
	Time     string
	Interval int32
	Site     []AlertStateSiteunit
}

func (this *AlertStateSite) FromArts(src *alertstate.GlobalResult) {
	this.Time = TimeFormUi(src.Time)
	this.Interval = alertstate.GLocalCach.WinWidth * 1000
	this.SiteFromArts(src.Site)
}
func (this *AlertStateSite) SiteFromArts(srcsite []alertstate.StateSiteUnit) {
	this.Site = make([]AlertStateSiteunit, 0)
	for _, src := range srcsite {
		unit := AlertStateSiteunit{}
		unit.FromArts(&src)
		this.Site = append(this.Site, unit)
	}
}
