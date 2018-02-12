package main

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

func containerReplicasByBarge(r *http.Request) *Response {

	query := `
SELECT "barge", SUM("replicas") AS "count"
FROM "Providers"
WHERE "replicas" > 0
  AND "barge" NOT LIKE '%test%'
GROUP BY "barge"
ORDER BY "count" DESC
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
