package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/PreetamJinka/catena"
	"github.com/VividCortex/siesta"
)
/*


type QueryDesc struct {
		Source string `json:"source"`
		Metric string `json:"metric"`
		Start  int64  `json:"start"`
		End    int64  `json:"end"`	
	}




	STEPS to parse a query:
	1) split query into its elements.
	2) build a list of descs for each type of query element
		1) sources: a new desc for each source
		2) metrics: a new desc for each metric
		3) start/end time:
			absolute: normal desc Start/End
			relative: calculate then place into Start/End
	3) return response

*/

type UserQuery struct {
		Sources []string `json:"sources"`
		Metrics []string `json:"metrics"`
		Start_Absolute  int64  `json:"start_absolute"`
		End_Absolute    int64  `json:"end_absolute"`	
	}

func query(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	db := c.Get(catenaKey).(*catena.DB)

	var params siesta.Params
	var user_query UserQuery //the user's query input
	
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&user_query)//set user_query




	downsample := params.Int64("downsample", 0, "A downsample value of averages N points at a time")
	err = params.Parse(r.Form)
	if err != nil {
		c.Set(errorKey, err.Error())
		return
	}

	var descs []catena.QueryDesc//the list of descs that will be used to call the db.
	
	dec = json.NewDecoder(r.Body)
	err = dec.Decode(&descs)

	//log.Println("-------------")
	//log.Println(descs)//[{himanshu mymetric  90 110}] 
	//NOTE:if no metric is provided then db will just store value without metric. weird.
	//todo: look into this.
	//log.Println("-------------")
	
	if err != nil {
		c.Set(errorKey, err.Error())
		return
	}


	now := time.Now().Unix()
	
	for i, desc := range descs { //i is index,  desc is querydesc.
		if desc.Start <= 0 {
			desc.Start += now
		}

		if desc.End <= 0 {
			desc.End += now
		}

		descs[i] = desc
	}


	/* db.Query takes in a list of 
	type QueryDesc struct {
		Source string `json:"source"`
		Metric string `json:"metric"`
		Start  int64  `json:"start"`
		End    int64  `json:"end"`	
	}

	so create a bunch of these for advanced queries.

	*/
	resp := db.Query(descs)

	if *downsample <= 1 {
		c.Set(responseKey, resp)
		return
	}

	for i, series := range resp.Series {
		pointIndex := 0
		seenPoints := 1
		currentPartition := series.Points[0].Timestamp / *downsample
		for j, p := range series.Points {
			if j == 0 {
				continue
			}

			if p.Timestamp / *downsample == currentPartition {
				series.Points[pointIndex].Value += p.Value
				seenPoints++
			} else {
				currentPartition = p.Timestamp / *downsample
				series.Points[pointIndex].Value /= float64(seenPoints)
				pointIndex++
				seenPoints = 1
				series.Points[pointIndex] = p
			}

			if j == len(series.Points) {
				series.Points[pointIndex].Value /= float64(seenPoints)
			}
		}

		series.Points = series.Points[:pointIndex]
		resp.Series[i] = series
	}

	c.Set(responseKey, resp)
}
