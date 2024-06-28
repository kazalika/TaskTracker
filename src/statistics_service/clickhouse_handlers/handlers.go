package clickhouse_handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var conn clickhouse.Conn
var statisticsNames map[string]struct{}

func InitConnection() error {
	// Init statistics names
	statisticsNames = map[string]struct{}{
		"likes": {},
		"views": {},
	}

	// Get environment variables
	address, ok := os.LookupEnv("CLICKHOUSE_ADDRESS")
	if !ok {
		return errors.New("clickhouse address is not configured (CLICKHOUSE_ADDRESS environment variable)")
	}
	database, ok := os.LookupEnv("CLICKHOUSE_DB")
	if !ok {
		return errors.New("clickhouse database is not configured (CLICKHOUSE_DB environment variable)")
	}
	username, ok := os.LookupEnv("CLICKHOUSE_USER")
	if !ok {
		return errors.New("clickhouse username is not configured (CLICKHOUSE_USER environment variable)")
	}
	password, ok := os.LookupEnv("CLICKHOUSE_PASSWORD")
	if !ok {
		return errors.New("clickhouse password is not configured (CLICKHOUSE_PASSWORD environment variable)")
	}

	// Open connection
	var err error
	conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{address},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		return err
	}

	v, err := conn.ServerVersion()
	if err == nil {
		log.Println("ClickHouse server version:", v)
	}
	return err
}

func CloseConnection() error {
	return conn.Close()
}

func getTaskNumbericStat(task_id int32, parameter string) (result uint64, err error) {
	if _, ok := statisticsNames[parameter]; !ok {
		err = fmt.Errorf("there is not statistic named `%s`", parameter)
		return 0, err
	}

	query := fmt.Sprintf(`
	SELECT
		COUNT(DISTINCT username)
	FROM %s
	WHERE
		task_id = %v;
	`, parameter, task_id)

	// Result should be 1 row
	err = conn.QueryRow(context.Background(), query).Scan(&result)

	if result == 0 {
		// Task was not created yet
		err = fmt.Errorf("task %v wasn't created yet", task_id)
		return
	}

	// 1 like and 1 view was sent to create empty statistics. We should decrement it
	result -= 1

	return
}

type Statistics map[string]any

type TaskWithStatistics struct {
	TaskID     int32      `json:"task_id"`
	Author     string     `json:"author"`
	Statistics Statistics `json:"statistics"`
}

func GetTaskStatistics(taskID int32) (statistics Statistics, err error) {
	statistics = make(map[string]any)
	for parameter := range statisticsNames {
		v, err := getTaskNumbericStat(taskID, parameter)
		if err != nil {
			err = fmt.Errorf("`getTaskStat` failed with error: %w", err)
			return nil, err
		}
		statistics[parameter] = v
	}
	return statistics, nil
}

func GetTopTasksByParameter(parameter string, topSize int) (res []TaskWithStatistics, err error) {
	if _, ok := statisticsNames[parameter]; !ok {
		err = fmt.Errorf("there is not statistic named `%s`", parameter)
		return
	}

	// Get `topSize` task ids with the most `stat`
	query := fmt.Sprintf(`
	SELECT
    	task_id,
		task_author
	FROM %s
	GROUP BY
		task_id,
		task_author
	ORDER BY
	COUNT(DISTINCT username) DESC
	LIMIT %v
	`, parameter, topSize)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return
	}

	// Iterate through result of query and get statistics
	for rows.Next() {
		var t TaskWithStatistics
		err = rows.Scan(&t.TaskID, &t.Author)
		if err != nil {
			return
		}

		t.Statistics, err = GetTaskStatistics(t.TaskID)
		if err != nil {
			return
		}

		res = append(res, t)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

type UserWithLikes struct {
	Author string `json:"author"`
	Likes  int64  `json:"likes"`
}

func GetTopUsersByLikes(usersCount int) ([]UserWithLikes, error) {
	// Get `usersCount` rows with most liked usernames and their like counts
	query := fmt.Sprintf(`
	SELECT
		task_author,
		COUNT() - COUNT(DISTINCT task_id) AS total_likes
	FROM likes
	GROUP BY task_author
	ORDER BY total_likes DESC
	LIMIT %v;
	`, usersCount)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var result []UserWithLikes

	// Iterate through result of query
	for rows.Next() {
		var t UserWithLikes
		err := rows.Scan(&t.Author, &t.Likes)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
