package main

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

//ShipmentEnvironmentsByBargeResult ...
type ShipmentEnvironmentsByBargeResult struct {
	Barge string `json:"barge"`
	Count int    `json:"count"`
}

func shipmentEnvironmentsByBarge(r *http.Request) *Response {

	db, err := sql.Open("postgres", os.Getenv("DB_CONNECT"))
	check(err)
	defer db.Close()

	sql := `
select
  "Providers".barge, count(*) as count
from
  "Providers",
  "Environments",
  "Shipments"
where
  "Providers"."environmentId" = "Environments".composite
  and "Environments"."shipmentId" = "Shipments"."name"
group by 
  "Providers".barge
order by 
  count(*) desc
;	
`
	var results []ShipmentEnvironmentsByBargeResult
	rows, err := db.Query(sql)
	check(err)
	defer rows.Close()
	for rows.Next() {
		var barge string
		var count int
		err := rows.Scan(&barge, &count)
		check(err)
		results = append(results, ShipmentEnvironmentsByBargeResult{
			Barge: barge,
			Count: count,
		})
	}
	err = rows.Err()
	check(err)

	return DataJSON(http.StatusOK, results, nil)
}
