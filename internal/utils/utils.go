package utils

import (
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/utils"
)

// GoAsync exec goroutine async
func GoAsync(msg string, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := utils.PrintStackTrace(msg, err)
				log.Errorf("GoAsync panic: %+v, %+v", err, stack)
			}
		}()
		f()
	}()
}
