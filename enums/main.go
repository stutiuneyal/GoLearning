package main

import "fmt"

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
)

func (l LogLevel) getLogLevel() string {
	switch l {
	case LevelTrace:
		return "Trace"
	case LevelDebug:
		return "Debug"
	default:
		return "Other"
	}
}

type ProbeTye int

const (
	LivenessProbe ProbeTye = iota
	StartupProbe
	ReadinessProbe
)

func main() {
	fmt.Println(LevelTrace.getLogLevel())
	fmt.Println(ReadinessProbe)
}
