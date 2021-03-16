//Package swagger This is a GPIO service implementation
package swagger

import (
	"net/http"
	"time"
)

//Logger Logs each request
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		lg.Debug(

			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
