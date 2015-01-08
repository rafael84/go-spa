package middleware

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

func WebLogger(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, req)
	res := rw.(negroni.ResponseWriter)
	defer func() {
		elapsed := time.Since(start)
		log.WithFields(log.Fields{
			"elapsed": elapsed,
			"method":  req.Method,
			"host":    req.URL.Host,
			"path":    req.URL.Path,
			"query":   req.URL.RawQuery,
			"status":  res.Status(),
			"size":    res.Size(),
		}).Info(req.Method + " " + req.URL.Path)
	}()
}
