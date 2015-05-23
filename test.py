#just using this for now. not actual tests

import requests
import json
import pprint



#old query
#curl -XGET -H 'Content-Type: application/json' -d 
#'[{"source": "himanshu", "metric": "mymetric", "start": 90, 
#"end": 110}]' localhost:8080/query



#insert
#curl -XPOST -H 'Content-Type: application/json' -d 
#'[{"source": "my.source", "metric": "my.metric", "timestamp": 100, "value": 0.0}]' 
#localhost:8080/insert
#{"data":null}



# type Row struct {
# 	Source    string  `json:"source"`
# 	Metric    string  `json:"metric"`
# 	Timestamp int64   `json:"timestamp"`
# 	Value     float64 `json:"value"`
# }

# type Rows []Row
data = []
for i in range(1,1000):
	row = {
				"source":"himanshu",
				"metric":"mymetric",
				"timestamp":i,
				"value":i
			}
	data.append(row) 

		
		
        	
	
#pprint.pprint(data)
data=json.dumps(data)
#r = requests.post('http://localhost:8080/insert', data= data)





data={
	"source_metric":[
			["himanshu","mymetric"],
			["my.source","my.metric"]
			],
	"start":"70+5s",
	"end":"80s"
	#"end": "2*10year-20year+20s*5s+20",
	#"limit":1
}
data=json.dumps(data)
r = requests.get('http://localhost:8080/query', data= data)
pprint.pprint(r.json())


