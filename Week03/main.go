package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func main() {
	// graceful shutdown
	stopCh := setUpSignalHandler()

	group, ctx := errgroup.WithContext(context.Background())

	// http server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        &handler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// pprof server
	pprof := &http.Server{Addr: ":6060", Handler: http.DefaultServeMux}

	// run http server
	group.Go(func() error {
		log.Println("server is running on 8080.")
		return server.ListenAndServe()
	})

	// enable pprof
	// http://localhost:6060/debug/pprof/
	group.Go(func() error {
		log.Println("pprof is running on 6060, address is http://localhost:6060/debug/pprof/")
		return pprof.ListenAndServe()
	})

	// shutdown server
	go func() {
		<-stopCh

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server shutdown failed: %v", err)
			return
		}
		log.Println("server shutdown gracefully.")
	}()

	// shutdown pprof
	go func() {
		<-stopCh

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := pprof.Shutdown(ctx); err != nil {
			log.Printf("pprof shutdown failed: %v", err)
			return
		}
		log.Println("pprof shutdown gracefully.")
	}()

	if err := group.Wait(); err != nil {
		log.Printf("errgroup: %v", err)
	}

	<-ctx.Done()
	log.Printf("app exit gracefully. %v", ctx.Err())
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()
	_, err = w.Write([]byte("hello"))
	if err != nil {
		return
	}
}

func setUpSignalHandler() <-chan struct{} {

	// ensure setUpSignalHandler only call once.
	close(onlyOneSignalHandler)

	stop := make(chan struct{})

	c := make(chan os.Signal, 2)
	// signal.Notify causes package signal to relay incoming signals to c.
	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will.
	signal.Notify(c, shutdownSignals...)

	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}
