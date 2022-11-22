package modbusRTU

import "log"

type RTUResponseFrame struct {
	TimeoutInMilliseconds int
	SlaveAddress          byte
	FunctionCode          byte
	StartRegister         uint16
	NumberOfRegisters     uint16
	Data                  []byte
}

type Register [2]byte

const (
	MODBUS_PORT         = 502
	RTU_FRAME_MAXSIZE   = 512
	ASCII_FRAME_MAXSIZE = 512
	TCP_FRAME_MAXSIZE   = 260
)

func GenerateResponse(slaveAddress, functionCode byte, data []byte) []byte {

	frame := new(RTUResponseFrame)
	frame.SlaveAddress = slaveAddress
	frame.FunctionCode = functionCode
	frame.Data = data

	// generate the ADU from the RTU frame
	return frame.GenerateRTUResponseFrame()
}

// GenerateRTUResponseFrame is a method corresponding to a RTUFrame object which
// returns a byte array representing the associated serial line/RTU
// application data unit (ADU)
func (frame *RTUResponseFrame) GenerateRTUResponseFrame() []byte {

	packetLen := 3
	dataLen := len(frame.Data)
	if dataLen > 0 {
		packetLen = RTU_FRAME_MAXSIZE
	}

	packet := make([]byte, packetLen)
	packet[0] = frame.SlaveAddress
	packet[1] = frame.FunctionCode
	packet[2] = 2
	bytesUsed := 3

	for i := 0; i < dataLen; i++ {
		packet[(bytesUsed + i)] = frame.Data[i]
	}
	bytesUsed += dataLen

	log.Printf("Before CRC: %x\n", packet[:bytesUsed])

	// add the crc to the end
	packet_crc := crc(packet[:bytesUsed])
	//packet[bytesUsed] = byte(packet_crc & 0xff)
	//packet[(bytesUsed + 1)] = byte(packet_crc >> 8)
	packet[bytesUsed] = byte(packet_crc & 0xff)
	packet[(bytesUsed + 1)] = byte(packet_crc >> 8)
	bytesUsed += 2
	log.Printf("After CRC: %x\n", packet[:bytesUsed])

	return packet[:bytesUsed]
}

// crc computes and returns a cyclic redundancy check of the given byte array
func crc(data []byte) uint16 {
	var crc16 uint16 = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		crc16 ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if crc16&0x0001 > 0 {
				crc16 = (crc16 >> 1) ^ 0xA001
			} else {
				crc16 >>= 1
			}
		}
	}
	return crc16
}
