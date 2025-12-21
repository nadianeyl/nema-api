package helper

import (
	"fmt"
	"sync"

	"github.com/nadianeyl/nema-api/internal/jsonlog"
)

func Background(logger *jsonlog.Logger, wg *sync.WaitGroup, fn func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		defer func() {
			if err := recover(); err != nil {
				logger.LogError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
