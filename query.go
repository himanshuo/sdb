package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"math/big"
	"strconv"
	"sort"


	"github.com/PreetamJinka/catena"
	"github.com/VividCortex/siesta"
	"github.com/himanshuo/evaler"

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
		SourceMetric [][]string `json:"source_metric"`
		Start  string  `json:"start"`
		End    string  `json:"end"`
		Limit int64 `json:"limit"`
		Descending bool `json:"descending"`
		Filter []QueryFilter `json:"ascending"`
		Aggregators []QueryAggregator `json:"aggregators"`
		MinDuration int64 `json:"min_duration"`
	}
type  QueryFilter struct{
	Metric string `json:"metric"`
	Operator string `json:"operator"`
	Value string `json:"value"`
}
type  QueryAggregator struct{
	Name string `json:"name"`
	Metric string `json:"metric"`
	//...add more fields as needed.
}
func query(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	db := c.Get(catenaKey).(*catena.DB)

	//var params siesta.Params
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
     	newDesc.Start = handleTime(userQuery.Start)
     	newDesc.End = handleTime(userQuery.End)

     	//add new desc to descs
     	descs = append(descs, *newDesc)

	}
	log.Println("---------descs-----------------")
	log.Println(descs)
	log.Println("-------------------------------")



	





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

	//NOTE: we are doing stuff like limit, ascending,... on POINTS of series. Not series themselves.
	for i, _ := range resp.Series{



		//handle ascending/descending. user specifies descending or not. if nothing given, then assume ascending.
		//always sort since it is time series.
		if userQuery.Descending {
			sort.Sort(ByDescendingTime{resp.Series[i].Points})
		} else {
			sort.Sort(ByAscendingTime{resp.Series[i].Points})
		}

		//handle duration
		//previous. trash.
		//a = append(    a[:i],     a[i+1:]...  )
		// if userQuery.Duration >0 && resp.Series[i].End - resp.Series[i].Start < userQuery.Duration {
		// 	resp.Series = append(resp.Series[0:i-1],resp.Series[i+1:] )//remov
		// 	// i=i-1
		// }
		//
		//new. correct sudocode
		// point_to_insert = cur_serious[0]
		
		// for all points in cur_series:
		// 	if curpoint.timestamp - last inserted point < DURATION:
		// 		do nothing 
		// 	else:
		// 		retpoints.add(point_to_insert)
		// 		point_to_insert = curpoint
		
		if (len(resp.Series[i].Points) > 0) && (userQuery.Duration > 0) {
			var newPoints []catena.Point
			point_to_insert := resp.Series[i].Points[0]
			for _, cur_point := range resp.Series[i].Points {
			    if (cur_point.Timestamp - point_to_insert.Timestamp) >= userQuery.Duration {
			    	newPoints = append(newPoints, point_to_insert)
			    	point_to_insert = cur_point
			    }
			}

			resp.Series[i].Points = newPoints
			//resp.Start and resp.End also change
			resp.Series[i].Start = newPoints[0].timestamp
			resp.Series[i].End = newPoints[len(newPoints)-1].timestamp
		}


		//handle limit
		if userQuery.Limit >0 {			
			resp.Series[i].Points = resp.Series[i].Points[:userQuery.Limit]	 
			resp.Series[i].End = resp.Series[i].Points[len(resp.Series[i].Points)-1].timestamp
		}


			
		



		 
	}

	

	


	/*
	A QueryResponse is returned after querying the DB with a QueryDesc.
     type QueryResponse struct {
	     Series []Series `json:"series"`
     }

    A Series is an ordered set of points
	for a source and metric over a range
	of time.
	type Series struct {
		// First timestamp
		Start int64 `json:"start"`

		// Last timestamp
		End int64 `json:"end"`

		Source string `json:"source"`
		Metric string `json:"metric"`

		Points []Point `json:"points"`
	}

	*/


	log.Println("---------resp---------------")
	log.Println(resp)
	log.Println(resp.Series)
	log.Println("-----------------------------")





	c.Set(responseKey, resp)
}








func handleTime(relativeTime string) int64{
	/*
		relative time would look something like:
			13423423432-2d+5h
			now()-2d+5h
		 OR
		 	yyyy-mm-dd hh:mm:ss
	*/

	//if in yyyy-mm-dd hh:mm:ss format:
	//var dateRegex = regexp.MustCompile(`(\d{4}-\d{2}-\d{4} \d{2}:\d{2}:\d{2})`)
	//dateRegex.FindAllString(relativeTime ,-1)
	const customForm = "2006-01-02 15:04:00"
    t, err := time.Parse(customForm, relativeTime)
    var intTime int64
    if err!=nil{
    	//else just handle it regularly.
		var bigTime,err = evaler.Eval(relativeTime)
		if err!=nil{//TODO:handle errors properly.
			log.Println("malformed time")
			//c.Set(errorKey, err.Error())
			return 1
		}
		intTime, err = BigratToInt(bigTime)
		if err!=nil{ //TODO:handle errors properly.
			log.Println("malformed time")
			//c.Set(errorKey, err.Error())
		return 1
		}
    } else {

    	intTime = t.Unix()
    }

    

	

	return handleNegativeTime(intTime) 

}

func handleNegativeTime(possibleNegTime int64) int64 {
	now := time.Now().Unix()
		if possibleNegTime <= 0 { 
			//historical query. given negative number, you want some time that
			//is say 10 seconds before now(). handles that.
			return possibleNegTime + now
		}
		return possibleNegTime
}

// BigratToInt converts a *big.Rat to an int64 (with truncation);
// returns an error for integer overflows
func BigratToInt(bigrat *big.Rat) (int64, error) {
	float_string := bigrat.FloatString(0)
	return strconv.ParseInt(float_string, 10, 64)
}




//[]catena.Point
 // type Point struct {
 // 	Timestamp int64   `json:"timestamp"`
 // 	Value     float64 `json:"value"`
 // }
type Points []catena.Point

func (s Points) Len() int      { return len(s) }
func (s Points) Swap(i, j int) { s[i], s[j] = s[j], s[i] }//special syntax. CHECK IT OUT.

//implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded Organs value.
type ByAscendingTime struct{ Points }
func (p ByAscendingTime) Less(i, j int) bool { return p.Points[i].Timestamp < p.Points[j].Timestamp }



type ByDescendingTime struct{ Points }
func (p ByDescendingTime) Less(i, j int) bool { return p.Points[i].Timestamp > p.Points[j].Timestamp }

//usage
//sort.Sort(ByAscendingTime{s})