
CREATE TABLE IF NOT EXISTS views_queue (
  username String,
  task_id Int32,
  task_author String
) ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka:9092',
       kafka_topic_list = 'views',
       kafka_group_name = 'group1',
       kafka_format = 'JSONEachRow';


CREATE TABLE IF NOT EXISTS likes_queue (
  username String,
  task_id Int32,
  task_author String
) ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka:9092',
       kafka_topic_list = 'likes',
       kafka_group_name = 'group1',
       kafka_format = 'JSONEachRow';

CREATE TABLE IF NOT EXISTS views (
  username String,
  task_id Int32,
  task_author String
) ENGINE = ReplacingMergeTree()
ORDER BY (task_id, username);

CREATE TABLE IF NOT EXISTS likes (
  username String,
  task_id Int32,
  task_author String
) ENGINE = ReplacingMergeTree()
ORDER BY (task_id, username);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_views TO views AS
SELECT 
  username,
  task_id,
  task_author
FROM views_queue;

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_likes TO likes AS
SELECT 
  username,
  task_id,
  task_author
FROM likes_queue;
