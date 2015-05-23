sdb
===
sdb is a simple time series database.

Usage
---

Installation
1. `go get ./...`
2. `go build`
3. `./sdb -listen=:8080 -data=/tmp/sdb`

Insert
curl -XPOST -H 'Content-Type: application/json' -d '[{"source": "my.source", "metric": "my.metric", "timestamp": 100, "value": 0.0}]' localhost:8080/insert
{"data":null}
 
View Sources
  ∂ [15:36:10] [~]: curl localhost:8080/sources?start=100
{"data":["my.source"]}
 
Query
  ∂ [15:37:37] [~]: curl -XGET -H 'Content-Type: application/json' -d '[{"source": "my.source", "metric": "my.metric", "start": 90, "end": 110}]' localhost:8080/query
{"data":{"series":[{"start":100,"end":100,"source":"my.source","metric":"my.metric","points":[{"timestamp":100,"value":0}]}]}}






License
---
BSD
