package build

import "fmt"

var version string = "undefined"
var user string = "undefined"
var time string = "undefined"
var number string = "undefined"

// Version - complete version string
func Version() string {
	return fmt.Sprintf("%s compiled at %s by %s (build #%s)",
		version, time, user, number)
}

// VersionShort - short version string
func VersionShort() string {
	return version
}

// For AT+GMR

// ATVersion - AT version
func ATVersion() string {
	return "AT version:0.23.0.0(Apr 24 2015 21:11:01)"
}

// SDKVersion - SDK Version
func SDKVersion() string {
	return fmt.Sprintf("SDK version:%s(%s)", version, number)
}

// BinVersion - Bin version for AT
func BinVersion() string {
	return fmt.Sprintf("Bin version (%s %s): %s", user, number, version)
}

// CompileTime - Compile time string
func CompileTime() string {
	return fmt.Sprintf("compile time:%s", time)
}
