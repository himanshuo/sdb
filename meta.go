package main

import (
	"net/http"
	"time"

	"github.com/PreetamJinka/catena"
	"github.com/VividCortex/siesta"
)

func getSources(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	db := c.Get(catenaKey).(*catena.DB)

	var params siesta.Params
	start := params.Int64("start", -86400, "Starting timestamp")
	end := params.Int64("end", 0, "Ending timestamp")
	err := params.Parse(r.Form)
	if err != nil {
		c.Set(errorKey, err.Error())
		return
	}

	now := time.Now().Unix()

	if *start <= 0 {
		*start += now
	}

	if *end <= 0 {
		*end += now
	}

	c.Set(responseKey, db.Sources(*start, *end))
}

func getMetrics(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	db := c.Get(catenaKey).(*catena.DB)

	var params siesta.Params
	source := params.String("source", "", "Source")
	start := params.Int64("start", -86400, "Starting timestamp")
	end := params.Int64("end", 0, "Ending timestamp")
	err := params.Parse(r.Form)
	if err != nil {
		c.Set(errorKey, err.Error())
		return
	}

	if *source == "" {
		c.Set(errorKey, "source unspecified")
		return
	}

	now := time.Now().Unix()

	if *start <= 0 {
		*start += now
	}

	if *end <= 0 {
		*end += now
	}

	c.Set(responseKey, db.Metrics(*source, *start, *end))
}
