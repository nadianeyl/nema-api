package httputil

import "github.com/nadianeyl/nema-api/internal/jsonlog"

type HTTPUtil struct {
	Logger *jsonlog.Logger
}

func New(logger *jsonlog.Logger) HTTPUtil {
	return HTTPUtil{
		Logger: logger,
	}
}
