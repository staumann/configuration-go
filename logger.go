package config

import (
	"log"
)

//LogLevel for configuration plugin
type LogLevel int

const (
	//DEBUG LogLevel debug
	DEBUG LogLevel = 0
	//INFO LogLevel info
	INFO LogLevel = 1
	//WARN LogLevel warn
	WARN LogLevel = 2
	//ERROR LogLevel error
	ERROR LogLevel = 3
)

var (
	current LogLevel
)

func setLogLevel(level LogLevel) {
	current = level
}

func logMessage(level LogLevel, msg string) {
	if level >= current {
		log.Print(msg)
	}
}
