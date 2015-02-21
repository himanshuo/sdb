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
		SourceMetric [][]string `json:"sourceMetric"`
		Start  int64  `json:"start_absolute"`
		End    int64  `json:"end_absolute"`	
	}

func query(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	db := c.Get(catenaKey).(*catena.DB)

	var params siesta.Params
	var userQuery UserQuery //the user's query input
	


	//get the user input and deserialize it
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&userQuery)//set userQuery

	if err!=nil{
		log.Println("there was an error.")
		c.Set(errorKey, err.Error())
		return
	}
	log.Println("-----query info--------")
	log.Println(r.Body)
	log.Println(userQuery) //{[[himanshu mymetric]] 0 0}
	log.Println("-------------")




	//unclear what this is doing.
	downsample := params.Int64("downsample", 0, "A downsample value of averages N points at a time")
	err = params.Parse(r.Form)
	if err != nil {
		c.Set(errorKey, err.Error())
		return
	}

	var descs []catena.QueryDesc//the list of descs that will be used to call the db.
	
	//go through all the source_query pairs and put them into descs array.
	for _, sourceMetric := range userQuery.SourceMetric {
    	
    	//get the source and metric from the query.
     	source := sourceMetric[0]
     	metric := sourceMetric[1]
     	
     	//build desc
     	newDesc := new(catena.QueryDesc)
     	newDesc.Source = source
     	newDesc.Metric = metric
     	newDesc.Start = userQuery.Start
     	newDesc.End = userQuery.End

     	//add new desc to descs
     	descs = append(descs, *newDesc)

	}
	log.Println("---------descs-----------------")
	log.Println(descs)
	log.Println("-------------------------------")



	handleRelativeTime(descs)





	log.Println("---------descs-----------------")
	log.Println(descs)
	log.Println("-------------------------------")

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


	log.Println("----------resp-----------------")
	log.Println(resp)
	log.Println("-------------------------------")



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



func handleRelativeTime(descs []catena.QueryDesc){


	now := time.Now().Unix()
	
	for i, desc := range descs { //i is index,  desc is querydesc.
		if desc.Start <= 0 { 
			//historical query. given negative number, you want some time that
			//is say 10 seconds before now(). handles that.
			desc.Start += now
		}

		if desc.End <= 0 {
			//same as desc.Start.
			desc.End += now
		}
		descs[i] = desc//update the desc with neg value with this logic.
	}
}