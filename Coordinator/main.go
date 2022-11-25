package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"coordinator/api"
	"coordinator/db"
	"coordinator/reposistory"
	"coordinator/usecase"
)

var (
	port       = flag.Int("port", 8080, "port")
	batchSize  = flag.Int("batch-size", 1000, "batch size")
	numRecords = flag.Int("records", 100000, "total number of records")
	numKeys    = flag.Int("keys", 100, "total number of keys")
)

func main() {
	d, err := db.NewMySQL(10, 10)
	if err != nil {
		log.Fatal("could not connect to DB", err)
	}

	pr := reposistory.NewProgressRepo(d)
	pu := usecase.NewProgressUseCase(pr, *batchSize, *numRecords, *numKeys)
	s := api.NewServer(pu)

	log.Println("start listening on port ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), s))
}
