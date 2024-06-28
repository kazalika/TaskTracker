/*
 * 		-- Statistics Service --
 * - Returns statistics for tasks or users
 * - Holds connection to ClickHouse with statistic tables
 */

package main

import (
	"clickhouse_handlers"
	"log"
	"net/http"
	han "statistics_service/api_handlers"
)

func main() {
	router := han.NewRouter()
	err := clickhouse_handlers.InitConnection()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer clickhouse_handlers.CloseConnection()

	log.Printf("Statistics service is starting...")
	log.Fatal(http.ListenAndServe(":8090", router))
}
