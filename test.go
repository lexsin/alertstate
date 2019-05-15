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
	ss := make([]EntryRecord, 0)
	for {
		select {
		case <-ticker.C:
			count = 0
			send(ss)
			ss = make([]EntryRecord, 0)
		default:
			if count == 1 {
				//fmt.Println("count=10000")
				continue
			}
			sample := genSample()
			ss = append(ss, sample)
			count++
		}
	}
}

func send(ss []EntryRecord) {
	fmt.Println("send len = ", len(ss))
	for _, record := range ss {
		InputCh <- record
	}
}

func genSample() EntryRecord {
	rand.Seed(time.Now().UnixNano())
	sniffer := rand.Intn(9) + 1
	site := rand.Intn(9) + 1
	genre := rand.Intn(8) + 1

	return EntryRecord{
		Timestamp:   int64(time.Now().Unix()),
		Genre:       classIntStr[classType(genre)],
		Sniffer:     int64(sniffer),
		Sniffername: strconv.Itoa(sniffer),
		Site:        int32(site),
		Sitename:    strconv.Itoa(site),
	}
}
