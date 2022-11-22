package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dugwill/boiler/modbusRTU"
	"github.com/tarm/serial"
)

func main() {

	connection, err := ConnectRTU("COM6", 38400)
	if err != nil {
		log.Fatalln("Error opening port")
	}
	defer connection.Close()

	buf := make([]byte, 128)

	byteCnt, response, err := readCommand(connection, buf)
	if err != nil {
		log.Println("Error reading input", err)
	}

	log.Printf("Received %d bytes: %x\n", byteCnt, response)

	data := []byte{0x03, 0x36, 0x04, 0x36}

	adu := modbusRTU.GenerateResponse(0x01, 0x04, data)

	w, err := connection.Write(adu)
	if err != nil {
		log.Print("Error writing to Serial: ", err)
	}
	log.Printf("Wrote %d Bytes: %x", w, adu)
}

// ConnectRTU attempts to access the Serial Device for subsequent
// RTU writes and response reads from the modbus slave device
func ConnectRTU(serialDevice string, baudRate int) (io.ReadWriteCloser, error) {
	conf := &serial.Config{Name: serialDevice, Baud: baudRate, StopBits: 2}
	ctx, err := serial.OpenPort(conf)
	ctx.Flush()
	return ctx, err
}

// DisconnectRTU closes the underlying Serial Device connection
func DisconnectRTU(ctx io.ReadWriteCloser) {
	ctx.Close()
}

func readCommand(connection io.ReadWriteCloser, buf []byte) (byteCnt int, response []byte, err error) {

	var rBytes []byte

	for {
		n, err := connection.Read(buf)
		if err != nil {
			log.Println("Error reading from port: ", err)
			return 0, nil, err
		}

		log.Printf("First Read: %d bytes, %x\n", n, buf[:n])
		if n < 8 {
			byteCnt += n
			rBytes = make([]byte, n)
			copy(rBytes, buf)
			time.Sleep(100 * time.Millisecond)
			n, err = connection.Read(buf)
			if err != nil {
				log.Println("Error reading from port: ", err)
			}
			log.Printf("Second Read: %d bytes, %x\n", n, buf[:n])
			byteCnt += n
			rBytes = append(rBytes, buf[:n]...)
		}

		logMsg := fmt.Sprintf("Read %d bytes: ", byteCnt)
		for i := 0; i < byteCnt; i++ {
			logMsg = logMsg + fmt.Sprintf("%x ", rBytes[i])
		}
		log.Println(logMsg)

		if byteCnt > 8 {
			log.Printf("Too many bytes read: %d. Flushing buffer.", byteCnt)
			continue
		} else {
			break
		}
	}
	return
}
