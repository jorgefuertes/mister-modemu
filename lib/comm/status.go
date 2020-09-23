package comm

type st struct {
	st     int8
	cipmux bool
	echo   bool
}

var status st

func resetStatus() {
	status = st{
		st:     5,
		cipmux: false,
	}
	_, err := getOutboundIP()
	if err != nil {
		status.st = 5
	} else {
		status.st = 2
	}
}
