package logrusgce

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type httpEntry struct {
	Method string `json:"requestMethod"`
	Url    string `json:"requestUrl"`
	Agent  string `json:"userAgent"`
}

// PostprocessHttpRequest adds gce formatted logging of an http request
func PostprocessHttpRequest(fieldname string) PostProcessFunc {
	return func(logrusEntry *logrus.Entry, gceEntry logrus.Fields) {
		raw, ok := logrusEntry.Data[fieldname]
		if !ok {
			return
		}
		r, ok := raw.(*http.Request)
		if !ok {
			return // we want to error?
		}
		httpEntry := httpEntry{
			Method: r.Method,
			Url:    r.URL.EscapedPath(),
			Agent:  r.UserAgent(),
		}
		gceEntry["httpRequest"] = httpEntry
	}
}
