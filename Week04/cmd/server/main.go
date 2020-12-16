package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Promacanthus/Go-000/Week04/configs"

	"github.com/gorilla/mux"

	_ "time/tzdata"

	_ "github.com/go-sql-driver/mysql"
)

// @title Go Week04 HomeWork Example API
// @version 1.0 https://github.com/Promacanthus
// @description This is a sample server.

// @contact.name promacanthus
// @contact.url
// @contact.email promacanthus@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v1
func main() {
	cfg := configs.NewConfig()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	stopCh := setUpSignalHandler()
	group, _ := errgroup.WithContext(context.Background())

	db, err := sql.Open("mysql", cfg.User+":"+cfg.Password+"@tcp("+cfg.DataBaseAddress+")/"+cfg.TableName)
	if err != nil {
		log.Fatalf("failed connectiong to database: %v", err)
	}

	h := InitHandler(db)

	r := mux.NewRouter()
	r.HandleFunc("/v1/save/{name}", h.Save).Methods(http.MethodPost)
	r.HandleFunc("/v1/get/{name}", h.Get).Methods(http.MethodGet)
	r.HandleFunc("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:1323/swagger/doc.json"))).Methods(http.MethodGet)

	server := &http.Server{
		Addr:           ":80",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	group.Go(func() error {
		return server.ListenAndServe()
	})

	group.Go(func() error {
		<-stopCh
		return server.Shutdown(context.Background())
	})

	if err := group.Wait(); err != nil {
		log.Fatalf("errgroup: %v", err)
	}
}

func setUpSignalHandler() <-chan struct{} {
	stop := make(chan struct{})

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		close(stop)
		<-c
	}()

	return stop
}
