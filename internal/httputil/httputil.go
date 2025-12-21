package httputil

import "github.com/nadianeyl/nema-api/internal/jsonlog"

type HTTPUtil struct {
	logger *jsonlog.Logger
}

func New(l *jsonlog.Logger) HTTPUtil {
	return HTTPUtil{
		logger: l,
	}
}
