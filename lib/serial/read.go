func readFreesizedSerialPortLoop(port *serial.Port, readpipe chan []byte) error {
	readBuf := make([]byte, 256)
	var sumBuf = []byte{}

	readPointer := 0

	defer port.Close()

	for {
		num, err := port.Read(readBuf)
		if err == io.EOF {
			// No more data comes
			if readPointer > 0 {
				log.Debugf("read data to send: %v", sumBuf)
				readpipe <- sumBuf
				readPointer = 0
				sumBuf = []byte{}
				log.Debugf("sumBuf cleared: %v", sumBuf)
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("cannnot open serial port: serialPort: %v, Error: %v", port, err)
		}
		if num > 0 {
			readPointer += num
			for index := range readBuf[:num] {
				sumBuf = append(sumBuf, readBuf[index])
			}
			log.Debugf("read partial data: %v to sumBuf: %v", readBuf, sumBuf)
			continue
		}
	}
}