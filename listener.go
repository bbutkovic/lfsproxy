package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

type lfsPacket struct {
	dataType string
	data     interface{}
}

type outGaugeData struct {
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

type outSimData struct {
	Time    uint32
	AngVelX float32
	AngVelY float32
	AngVelZ float32
	Heading float32
	Pitch   float32
	Roll    float32
	AccelX  float32
	AccelY  float32
	AccelZ  float32
	VelX    float32
	VelY    float32
	VelZ    float32
	PosX    int32
	PosY    int32
	PosZ    int32
	ID      int32
}

func (lp *lfsPacket) toJSON() ([]byte, error) {
	var s []byte
	var err error
	switch d := lp.data.(type) {
	case outGaugeData:
		s, err = json.Marshal(d)
	case outSimData:
		s, err = json.Marshal(d)
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func outGaugeListener(address string, c chan lfsPacket) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	var decoded outGaugeData
	buffer := make([]byte, binary.Size(&decoded))

	for {
		if _, err := conn.Read(buffer); err != nil {
			log.Fatal(err)
		}

		if err := binary.Read(bytes.NewReader(buffer), binary.LittleEndian, &decoded); err != nil {
			log.Fatal(err)
		}

		packet := lfsPacket{
			dataType: "outGauge",
			data:     decoded,
		}

		select {
		case c <- packet:
		default:
		}
	}
}

func outSimListener(address string, c chan lfsPacket) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	var decoded outSimData
	buffer := make([]byte, binary.Size(&decoded))

	for {
		if _, err := conn.Read(buffer); err != nil {
			log.Fatal(err)
		}

		if err := binary.Read(bytes.NewReader(buffer), binary.LittleEndian, &decoded); err != nil {
			log.Fatal(err)
		}

		packet := lfsPacket{
			dataType: "outSim",
			data:     decoded,
		}

		select {
		case c <- packet:
		default:
		}
	}
}
