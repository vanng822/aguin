package utils

import (
	"log"
	"os"
	"fmt"
)

func GetLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.LstdFlags)
}