package app

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"net"

	"github.com/sigurn/crc16"
	"github.com/sinantomruk/miltron-launchsite/models"
)

func bytesToFloat32(in []byte) float32 {
	bits := binary.BigEndian.Uint32(in)
	res := math.Float32frombits(bits)
	return res
}

func checkCRC(data []byte) bool {
	table := crc16.MakeTable(crc16.CRC16_BUYPASS)
	crc := crc16.Checksum(data[:33], table)
	sum := binary.BigEndian.Uint16(data[33:35])
	return crc == sum
}

func ReadTelemetry(ctx context.Context, host string, port int, dataChan chan<- models.TelemetryData) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	buffer := make([]byte, 36)
	for {
		if _, err := conn.Read(buffer); err != nil {
			return err
		}

		if !checkCRC(buffer) {
			continue
		}

		dataChan <- models.TelemetryData{
			RocketID:     string(buffer[1:11]),
			Altitude:     bytesToFloat32(buffer[13:17]),
			Speed:        bytesToFloat32(buffer[17:21]),
			Acceleration: bytesToFloat32(buffer[21:25]),
			Thrust:       bytesToFloat32(buffer[25:29]),
			Temperature:  bytesToFloat32(buffer[29:33]),
		}
	}
}
