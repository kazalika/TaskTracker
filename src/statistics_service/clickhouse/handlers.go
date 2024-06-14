package clickhouse_handlers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var conn clickhouse.Conn

func InitConnection() error {
	var err error
	conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{"clickhouse:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
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

func GetTaskStat(task_id string, stat string) (result uint64, err error) {
	if stat != "likes" && stat != "views" {
		return 0, errors.New("there are only views and likes stats")
	}

	query := fmt.Sprintf(`
	SELECT
		COUNT(DISTINCT username)
	FROM %s
	WHERE
		task_id = '%s';
	`, stat, task_id)

	log.Printf("Get task %s for task with id = %s", stat, task_id)

	err = conn.QueryRow(context.Background(), query).Scan(&result)

	// result == 0 mean's that this task was not created yet (maybe it's impossible scenario)
	if result != 0 {
		// Decrement Admin's like or view that was send for creating empty statistics
		result -= 1
	}

	return
}

func GetTaskLikes(task_id string) (likes uint64, err error) {
	return GetTaskStat(task_id, "likes")
}

func GetTaskViews(task_id string) (likes uint64, err error) {
	return GetTaskStat(task_id, "views")
}

type TaskWithStatistics struct {
	TaskID string `json:"task_id"`
	Likes  uint64 `json:"likes"`
	Views  uint64 `json:"views"`
}

func GetTopTasksByStat(stat string, topSize int) ([]TaskWithStatistics, error) {
	if stat != "likes" && stat != "views" {
		return nil, errors.New("there are only views and likes stats")
	}
	query := fmt.Sprintf(`
	SELECT
    	task_id
	FROM %s
	GROUP BY
		task_id,
		task_author
	ORDER BY
	COUNT(DISTINCT username) DESC
	LIMIT %v
	`, stat, topSize)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var result []TaskWithStatistics

	for rows.Next() {
		var t TaskWithStatistics
		err := rows.Scan(&t.TaskID)
		if err != nil {
			return nil, err
		}

		// Get likes and views for this task
		t.Likes, err = GetTaskLikes(t.TaskID)
		if err != nil {
			return nil, err
		}
		t.Views, err = GetTaskViews(t.TaskID)
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

type UserWithStatistics struct {
	Author string `json:"author"`
	Likes  int64  `json:"likes"`
}

func GetTopUsersByLikes(usersCount int) ([]UserWithStatistics, error) {
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

	var result []UserWithStatistics

	for rows.Next() {
		var t UserWithStatistics
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
