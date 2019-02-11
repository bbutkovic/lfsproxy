package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {

}

type outGaugeStruct struct {
	Time        uint32
	Car         [4]byte // In C that would be a char[4]
	Flags       uint16  // In C that would be a WORD (2 bytes of data)
	Gear        byte
	PLID        byte
	Speed       float32
	RPM         float32
	Turbo       float32
	EngTemp     float32
	Fuel        float32
	OilPressure float32
	OilTemp     float32
	DashLights  uint32
	ShowLights  uint32
	Throttle    float32
	Brake       float32
	Clutch      float32
	Display1    [16]byte
	Display2    [16]byte
	ID          int32
}

func outGaugeListener(address string, c chan outGaugeStruct) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	var decoded outGaugeStruct
	buffer := make([]byte, binary.Size(&decoded))

	for {
		if _, err := conn.Read(buffer); err != nil {
			log.Fatal(err)
		}

		if err := binary.Read(bytes.NewReader(buffer), binary.LittleEndian, &decoded); err != nil {
			log.Fatal(err)
		}

		select {
		case c <- decoded:
			fmt.Println("Wrote!")
		default:
			fmt.Print(".")
		}
	}
}
