package main

import (
	"fmt"
	"io"
	"log"

	"github.com/dugwill/modbus"
)

func main() {

	//port := `/dev/ttyUSB1`
	port := `com6`
	recChan := make(chan []byte)

	//t := triangleTube.NewBoiler(1)

	var rtu modbus.Connection = &modbus.RTU{SerialDevice: port, BaudRate: 38400, StopBits: 1}

	err := rtu.Connect()
	if err != nil {
		log.Fatalln("Error opening serial port: ", err)
	} else {
		log.Println("Serial connetion successful")
	}
	defer disconnect(rtu)

	go rtu.Read(recChan)

	for {
		response := <-recChan
		log.Printf("Received %d bytes: %x\n", len(response), response)
	}
}

func readCommand(connection io.ReadWriteCloser, buf []byte) (byteCnt int, response []byte, err error) {

	var rBytes []byte

	for {

		logMsg := fmt.Sprintf("Read %d bytes test: ", byteCnt)
		for i := 0; i < byteCnt; i++ {
			logMsg = logMsg + fmt.Sprintf("%x ", rBytes[i])
		}
		log.Println(logMsg)

		if byteCnt > 8 {
			log.Printf("Too many bytes read: %d. Flushing buffer.", byteCnt)
			continue
		} else {
			//validateChecksum
		}
	}
	// return buf without checksum
	return byteCnt, buf[:6], nil
}

func disconnect(rtu modbus.Connection) {
	log.Println("shutting down.")
	rtu.Disconnect()

}
