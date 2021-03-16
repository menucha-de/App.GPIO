/*
 * GPIO
 *
 * API version: 1.0.0
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	gpio "github.com/menucha-de/App.GPIO/gpio"
	capture "github.com/menucha-de/capture"
	transport "github.com/menucha-de/transport"
	"github.com/menucha-de/logging"
	"github.com/menucha-de/utils"
)

var log *logging.Logger = logging.GetLogger("gpio")

func call(id int, state int) {
	m := capture.CaptureData{Date: time.Now(), Device: "gpio", Field: strconv.Itoa(id), Value: strconv.Itoa(state)}
	topic := strconv.Itoa(id)
	capture.Pub.Publish(topic, m)
	gpio.HandleMessages1(topic, m)
	gpio.Clients.Mu.Lock()
	defer gpio.Clients.Mu.Unlock()
	if !capture.HasEnabledReports() && len(gpio.Clients.Clients) == 0 {
		gpio.SetEnable(id, false)
	}

}

func main() {
	var port = flag.Int("p", 8081, "port")
	flag.Parse()

	gpio.SetCallback(call)
	gpio.AddRoutes(logging.LogRoutes) //must be before router initialization
	gpio.AddRoutes(transport.TransportRoutes)

	router := gpio.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(notFound)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	srv1 := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	done := make(chan os.Signal, 1)
	errs := make(chan error)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}
	}()
	go func() {
		if err := srv1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}

	}()
	log.Info("Server Started")
	select {
	case err := <-errs:
		log.Error("Could not start serving service due to (error: %s)", err)
	case <-done:
		log.Info("Server Stopped")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		logging.Close()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Failed:", err.Error())
	}
	if err := srv1.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Failed:", err.Error())
	}
	log.Info("Server Exited Properly")

}
func notFound(w http.ResponseWriter, r *http.Request) {
	if !(r.Method == "GET") {
		w.WriteHeader(404)
		return
	}

	file := "./www" + html.EscapeString(r.URL.Path)
	if file == "./www/" {
		file = "./www/index.html"
	}
	if file == "./www/openapi/" {
		file = "./www/openapi/index.html"
	}

	if utils.FileExists(file) {
		http.ServeFile(w, r, file)
	} else {
		w.WriteHeader(404)
	}

}
