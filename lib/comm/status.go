package comm

type connection struct {
	t    string
	ip   string
	port int16
	keep int16
}

type st struct {
	st          int8
	cipmux      int8
	echo        bool
	connections [5]*connection
}

var status st

func resetStatus() {
	status = st{
		st:     5,
		cipmux: 0,
		echo:   true,
	}
	_, err := getOutboundIP()
	if err != nil {
		status.st = 5
	} else {
		status.st = 2
	}
}
