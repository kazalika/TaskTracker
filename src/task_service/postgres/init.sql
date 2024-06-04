CREATE TABLE IF NOT EXISTS task_service_db (
    id SERIAL PRIMARY KEY,
    creator_username TEXT NOT NULL,
    task_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL
);