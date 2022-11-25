package main

import (
	"flag"
	"log"
	"net/http"

	"Worker/constant"
	"Worker/db"
	"Worker/reposistory"
	"Worker/service"
	"Worker/usecase"
)

var numWorker = flag.Int("worker", 3, "number of workers")

func main() {
	log.Println("Worker starting...")

	d, err := db.NewMySQL(10, 10)
	if err != nil {
		log.Fatal("could not connect to DB", err)
	}

	kr := reposistory.NewKeyRepo(d)
	rr := reposistory.NewRecordRepo(d)
	sr := reposistory.NewSignatureRepo(d)
	ps := service.NewProgressService(&http.Client{Timeout: constant.APITimeout})

	for i := 0; i < *numWorker; i++ {
		w := usecase.NewWorker(kr, rr, sr, ps)
		go w.Start()
	}

	exit := make(chan struct{})
	<-exit
}
