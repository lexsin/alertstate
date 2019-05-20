package alertstate

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func StartTest() {
	ticker := time.NewTicker(1 * time.Second)
	count := 0
	tickercount := 0
	ss := make([]EntryRecord, 0)
	for {
		select {
		case <-ticker.C:
			//			if tickercount == 2 {
			//				return
			//			}
			tickercount++
			count = 0
			send(ss)
			fmt.Println("len ss=", len(ss))
			ss = make([]EntryRecord, 0)
		default:
			if count == 100000 {
				fmt.Println("count=10000")
				continue
			}
			sample := genSample(time.Now().Unix())
			ss = append(ss, sample)
			count++
		}
	}
}

func send(ss []EntryRecord) {
	fmt.Println("send len = ", len(ss))
	count := 0
	for _, record := range ss {
		count++
		InputCh <- record
	}
	fmt.Println("send final count=", count)
}

func genSample(t int64) EntryRecord {
	rand.Seed(time.Now().UnixNano())
	sniffer := rand.Intn(9) + 1
	site := rand.Intn(9) + 1
	genre := rand.Intn(8) + 1

	return EntryRecord{
		Timestamp:   int64(t),
		Genre:       classIntStr[classType(genre)],
		Sniffer:     int64(sniffer),
		Sniffername: strconv.Itoa(sniffer),
		Site:        int32(site),
		Sitename:    strconv.Itoa(site),
	}
}
