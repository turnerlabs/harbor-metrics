package main

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

//ShipmentEnvironmentsByBargeResult ...
type ShipmentEnvironmentsByBargeResult struct {
	Barge string `json:"barge"`
	Count int    `json:"count"`
}

func shipmentEnvironmentsByBarge(r *http.Request) *Response {

	query := `
select
  "Providers".barge, count(*) as count
from
  "Providers",
  "Environments",
  "Shipments"
where
  "Providers"."environmentId" = "Environments".composite
	and "Environments"."shipmentId" = "Shipments"."name"
	and "Providers".barge NOT LIKE '%test%'
group by 
  "Providers".barge
order by 
  count(*) desc
;	
`
	var results []ShipmentEnvironmentsByBargeResult
	err := dbQuery(query, func(rows *sql.Rows) {
		var barge string
		var count int
		err := rows.Scan(&barge, &count)
		check(err)
		results = append(results, ShipmentEnvironmentsByBargeResult{
			Barge: barge,
			Count: count,
		})
	})
	check(err)

	return DataJSON(http.StatusOK, results, nil)
}
