package helper

import (
	"fmt"

	"github.com/nadianeyl/nema-api/internal/jsonlog"
)

func Background(logger *jsonlog.Logger, fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
