/*
 * 		-- Statistics Service --
 * Returns statistics for tasks or users:
 *		1. likes and views
 *		2. top tasks by views or likes (parameter)
 * 		3. top users by likes
 *
 * Holds connection to ClickHouse with statistic tables
 */

package main

import (
	"clickhouse_handlers"
	"log"
	"net/http"
	sw "statistics_service/go"
)

func main() {
	router := sw.NewRouter()
	err := clickhouse_handlers.InitConnection()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer clickhouse_handlers.CloseConnection()

	log.Printf("Statistics service is starting...")
	log.Fatal(http.ListenAndServe(":8090", router))
}
