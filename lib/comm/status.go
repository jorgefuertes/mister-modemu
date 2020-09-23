package comm

type st struct {
	st     int8
	cipmux bool
}

var status st = st{
	st:     5,
	cipmux: false,
}

func checkStatus() {
	_, err := getOutboundIP()
	if err != nil {
		status.st = 5
	} else {
		status.st = 2
	}
}
