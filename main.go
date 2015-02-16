package main

import (
	"flag"
	"log"
	"net/http"
	
	"fmt"

	"github.com/PreetamJinka/catena"
	"github.com/VividCortex/siesta"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "Listen address")
	dataDir := flag.String("data", "/opt/sdb", "Data directory")
	flag.Parse()

	db, err := catena.OpenDB(*dataDir)
	if err != nil {
		log.Fatal(err)
	}

	service := siesta.NewService("/")
	service.AddPre(func(c siesta.Context, w http.ResponseWriter, r *http.Request) {
		c.Set(catenaKey, db)
	})

	service.AddPost(jsonMarshaler)

	service.Route("GET", "/sources", "Gets a list of sources", getSources)
	service.Route("GET", "/source/:source/metrics", "Gets a list of metrics for a source", getMetrics)

	service.Route("POST", "/insert", "Inserts rows into the database", insert)
	service.Route("GET", "/query", "Queries the database", query)

	fmt.Println(service)

	log.Fatal(http.ListenAndServe(*listenAddr, service))
}
