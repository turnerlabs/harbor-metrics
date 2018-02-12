package main

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

//ShipmentEnvironmentsByGroupResult ...
type ShipmentEnvironmentsByGroupResult struct {
	Group string `json:"group"`
	Count int    `json:"count"`
}

func shipmentEnvironmentsByGroup(r *http.Request) *Response {

	query := `
select
  "Shipments".group, count(*) as count
from
  "Providers",
  "Environments",
  "Shipments"
where
  "Providers"."environmentId" = "Environments".composite
	and "Environments"."shipmentId" = "Shipments"."name"
	and "Providers".barge NOT LIKE '%test%'
group by 
	"Shipments".group
having
  count(*) > 45
order by 
  count(*) desc
;	
`
	var results []ShipmentEnvironmentsByGroupResult
	err := dbQuery(query, func(rows *sql.Rows) {
		var group string
		var count int
		err := rows.Scan(&group, &count)
		check(err)
		results = append(results, ShipmentEnvironmentsByGroupResult{
			Group: group,
			Count: count,
		})
	})
	check(err)

	return DataJSON(http.StatusOK, results, nil)
}
