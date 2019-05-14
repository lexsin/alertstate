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
	ss := make([]entryRecord, 0)
	for {
		select {
		case <-ticker.C:
			count = 0
			send(ss)
			ss = make([]entryRecord, 0)
		default:
			if count == 50000 {
				//fmt.Println("count=10000")
				continue
			}
			sample := genSample()
			ss = append(ss, sample)
			count++
		}
	}
}

func send(ss []entryRecord) {
	fmt.Println("send len = ", len(ss))
	for _, record := range ss {
		InputCh <- record
	}
}

func genSample() entryRecord {
	rand.Seed(time.Now().UnixNano())
	sniffer := rand.Intn(9) + 1
	site := rand.Intn(9) + 1
	class := rand.Intn(8) + 1

	return entryRecord{
		timestamp:   int64(time.Now().Unix()),
		class:       classType(class),
		sniffer:     int32(sniffer),
		sniffername: strconv.Itoa(sniffer),
		site:        int32(site),
		sitename:    strconv.Itoa(site),
	}
}
