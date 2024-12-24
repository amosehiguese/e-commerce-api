package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func Get() *zap.Logger {
	once.Do(func() {
		var err error
		instance, err = zap.NewProduction()
		if err != nil {
			log.Fatal(err)
		}
	})
	return instance
}
