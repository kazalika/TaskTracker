package task_service

import (
	"context"
	"database/sql"
	"sync/atomic"

	postgres "postgres"
	task_servicepb "task_service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	task_servicepb.UnimplementedTaskServiceServer
	db          *sql.DB
	taskCounter atomic.Int64
}

func NewServer() (server *Server, err error) {
	server = &Server{}
	server.taskCounter.Store(0)
	server.db = postgres.InitPostgreSQLClient()
	return
}

func GenerateTaskID(s *Server) int32 {
	return int32(s.taskCounter.Add(1))
}

func (s *Server) CreateTask(ctx context.Context, request *task_servicepb.TaskContent) (*task_servicepb.TaskID, error) {
	taskID := GenerateTaskID(s)
	_, err := s.db.ExecContext(ctx, "INSERT INTO task_service_db (creator_username, task_id, title, description, status) VALUES ($1, $2, $3, $4, $5)", request.CreatorUsername, taskID, request.Title, request.Description, request.Status)
	if err != nil {
		return &task_servicepb.TaskID{Id: taskID}, status.Errorf(codes.Internal, "insert into db has been failed, taskID: %v", taskID)
	}

	return &task_servicepb.TaskID{Id: taskID}, nil
}
func (s *Server) UpdateTask(ctx context.Context, request *task_servicepb.Task) (*task_servicepb.TaskID, error) {
	txn, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return &task_servicepb.TaskID{}, status.Errorf(codes.Internal, "failed to start transaction")
	}
	defer txn.Rollback()

	count := 0
	txn.QueryRowContext(ctx, "SELECT COUNT(*) FROM task_service_db WHERE creator_username = $1 AND task_id = $2", request.Task.CreatorUsername, request.Id).Scan(&count)
	if count == 0 {
		return &task_servicepb.TaskID{Id: request.Id}, status.Errorf(codes.NotFound, "there is no task with ID=%v by user=%s", request.Id, request.Task.CreatorUsername)
	}

	_, err = txn.ExecContext(ctx, "UPDATE task_service_db SET title = $1, description = $2, status = $3 WHERE creator_username = $4 AND task_id = $5", request.Task.Title, request.Task.Description, request.Task.Status, request.Task.CreatorUsername, request.Id)
	if err != nil {
		return &task_servicepb.TaskID{Id: request.Id}, status.Errorf(codes.Internal, "failed to update task with ID=%v from user=%s", request.Id, request.Task.CreatorUsername)
	}

	err = txn.Commit()
	if err != nil {
		return &task_servicepb.TaskID{Id: request.Id}, status.Errorf(codes.Internal, "failed to commit transaction")
	}

	return &task_servicepb.TaskID{Id: request.Id}, nil
}

func (s *Server) DeleteTask(ctx context.Context, request *task_servicepb.RequestByID) (*task_servicepb.TaskID, error) {
	txn, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return &task_servicepb.TaskID{}, status.Errorf(codes.Internal, "failed to start transaction")
	}
	defer txn.Rollback()

	count := 0
	txn.QueryRowContext(ctx, "SELECT COUNT(*) FROM task_service_db WHERE creator_username = $1 AND task_id = $2", request.RequestorUsername, request.Id).Scan(&count)
	if count == 0 {
		return &task_servicepb.TaskID{Id: request.Id}, status.Errorf(codes.NotFound, "there is no task with ID=%v by user=%s", request.Id, request.RequestorUsername)
	}

	_, err = txn.ExecContext(ctx, "DELETE FROM task_service_db WHERE creator_username = $1 AND task_id = $2", request.RequestorUsername, request.Id)
	if err != nil {
		return &task_servicepb.TaskID{Id: request.Id}, status.Errorf(codes.Internal, "failed to delete task with ID=%v from user%s", request.Id, request.RequestorUsername)
	}

	err = txn.Commit()
	if err != nil {
		return &task_servicepb.TaskID{Id: request.Id}, status.Errorf(codes.Internal, "failed to commit transaction")
	}

	return &task_servicepb.TaskID{Id: request.Id}, nil
}

func (s *Server) GetTaskById(ctx context.Context, request *task_servicepb.RequestByID) (*task_servicepb.Task, error) {
	var title, description, taskStatus, creator string
	err := s.db.QueryRowContext(ctx, "SELECT title, description, status, creator_username FROM task_service_db WHERE task_id = $1", request.Id).Scan(&title, &description, &taskStatus, &creator)
	if err != nil {
		if err == sql.ErrNoRows {
			return &task_servicepb.Task{}, status.Errorf(codes.NotFound, "there is no task with ID=%v", request.Id)
		} else {
			return &task_servicepb.Task{}, status.Errorf(codes.Internal, "failed to get task with ID=%v", request.Id)
		}
	}

	return &task_servicepb.Task{
		Id: request.Id,
		Task: &task_servicepb.TaskContent{
			Title:           title,
			Description:     description,
			Status:          taskStatus,
			CreatorUsername: creator,
		},
	}, nil
}

func (s *Server) GetTaskList(ctx context.Context, request *task_servicepb.TaskPageRequest) (*task_servicepb.TaskList, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT task_id, title, description, status, creator_username FROM task_service_db ORDER BY task_id LIMIT $1 OFFSET $2", request.PageSize, request.Offset)
	if err != nil {
		return &task_servicepb.TaskList{}, status.Errorf(codes.Internal, "failed to get task with offset=%v, page_size=%v", request.Offset, request.PageSize)
	}
	defer rows.Close()

	var response task_servicepb.TaskList
	tasks_list := response.GetTasks()

	for rows.Next() {
		task := &task_servicepb.Task{
			Id:   0,
			Task: &task_servicepb.TaskContent{},
		}
		err = rows.Scan(&task.Id, &task.Task.Title, &task.Task.Description, &task.Task.Status, &task.Task.CreatorUsername)

		if err != nil {
			return &task_servicepb.TaskList{}, status.Errorf(codes.Internal, "failed to read field in a row")
		}

		tasks_list = append(tasks_list, task)
	}

	return &task_servicepb.TaskList{Tasks: tasks_list, PageSize: int32(len(tasks_list))}, nil
}
