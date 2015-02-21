#just using this for now. not actual tests

import requests
import json


#curl -XGET -H 'Content-Type: application/json' -d 
#'[{"source": "himanshu", "metric": "mymetric", "start": 90, 
#"end": 110}]' localhost:8080/query


data={
	"source_metric":[
			["himanshu","mymetric"]
			],
	"start":90,
	"end": 110
}
data=json.dumps(data)
r = requests.get('http://localhost:8080/query', data= data)
print r.json()


