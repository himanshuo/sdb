package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PreetamJinka/catena"
	"github.com/VividCortex/siesta"
)

func insert(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	db := c.Get(catenaKey).(*catena.DB)

	rows := catena.Rows{}


	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&rows)

	log.Println(rows)
	if err != nil {
		c.Set(errorKey, err.Error())
		return
	}
	

	err = db.InsertRows(rows)
	if err != nil {
		c.Set(errorKey, err.Error())
	}
}
