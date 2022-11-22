package modbusASCII

import (
	"bytes"
	"encoding/hex"
	"log"
)

type ASCIIResponseFrame struct {
	SlaveAddress byte
	FunctionCode byte
	Data         []byte
}

type Register [2]byte

const (
	MODBUS_PORT         = 502
	RTU_FRAME_MAXSIZE   = 512
	ASCII_FRAME_MAXSIZE = 512
	TCP_FRAME_MAXSIZE   = 260
)

func GenerateResponse(slaveAddress, functionCode byte, data []byte) []byte {

	frame := new(ASCIIResponseFrame)
	frame.SlaveAddress = slaveAddress
	frame.FunctionCode = functionCode
	frame.Data = []byte{0x3, 0x36}

	return frame.GenerateASCIIResponseFrame()
}

// GenerateASCIIFrame is a method corresponding to a ASCIIFrame object which
// returns a byte array representing the associated serial line/ASCII
// application data unit (ADU)
func (frame *ASCIIResponseFrame) GenerateASCIIResponseFrame() []byte {

	packetLen := 3
	if len(frame.Data) > 0 {
		packetLen += len(frame.Data) //+ 1
		if packetLen > ASCII_FRAME_MAXSIZE {
			packetLen = ASCII_FRAME_MAXSIZE
		}
	}

	packet := make([]byte, packetLen)
	packet[0] = frame.SlaveAddress
	packet[1] = frame.FunctionCode
	packet[2] = byte(0x02)
	bytesUsed := 3

	for i := 0; i < len(frame.Data); i++ {
		packet = append(packet, frame.Data[i])
	}
	bytesUsed += len(frame.Data)

	log.Printf("BytesUsed: %d", bytesUsed)

	// add the lrc to the end
	packet_lrc := lrc(packet[:bytesUsed])
	packet[bytesUsed] = byte(packet_lrc)
	bytesUsed += 1

	log.Printf("BytesUsed: %d", bytesUsed)

	// Convert raw bytes to ASCII packet
	ascii_packet := make([]byte, bytesUsed*2+3)
	hex.Encode(ascii_packet[1:], packet)

	asciiBytesUsed := bytesUsed*2 + 1

	// Frame the packet
	ascii_packet[0] = ':'                 // 0x3A
	ascii_packet[asciiBytesUsed] = '\r'   // CR 0x0D
	ascii_packet[asciiBytesUsed+1] = '\n' // LF 0x0A
	asciiBytesUsed += 2

	return bytes.ToUpper(ascii_packet[:asciiBytesUsed])
}

// Modbus ASCII does not use CRC, but Longitudinal Redundancy Check.
// lrc computes and returns the 2's compliment (-) of the sum of the given byte
// array modulo 256
func lrc(data []byte) uint8 {
	var sum uint8 = 0
	var lrc8 uint8 = 0
	for _, b := range data {
		sum += b
	}
	lrc8 = uint8(-int8(sum))
	return lrc8
}
