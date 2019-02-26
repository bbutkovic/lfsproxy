package main

import (
	"fmt"
	"log"
)

func main() {
	c := make(chan lfsPacket)
	go outGaugeListener("127.0.0.1:2022", c)
	chans := makeStructChannels(3)
	go crossbar(c, nil, chans)

	for {
		data := <-chans[0]

		switch d := data.data.(type) {
		case outGaugeData:
			fmt.Println(d.RPM)
		}
	}
}

func makeJSONChannels(n int) []chan []byte {
	chans := make([]chan []byte, n)
	for i := range chans {
		chans[i] = make(chan []byte)
	}
	return chans
}

func makeStructChannels(n int) []chan lfsPacket {
	chans := make([]chan lfsPacket, n)
	for i := range chans {
		chans[i] = make(chan lfsPacket)
	}
	return chans
}

func crossbar(in chan lfsPacket, outJSON []chan []byte, outStruct []chan lfsPacket) {
	for {
		select {
		case data := <-in:
			for _, c := range outStruct {
				select {
				case c <- data:
				default:
				}
			}
			if len(outJSON) > 0 {
				s, err := data.toJSON()
				if err != nil {
					log.Panic(err)
				}
				for _, c := range outJSON {
					select {
					case c <- s:
					default:
					}
				}
			}
		}
	}
}
