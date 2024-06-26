## Как запустить?

```bash
docker compose build
docker compose up
```

## Примечания про auth_service

1. Основной сервис (авторизация) написан на Go с использованием кодогенерации Swagger.

2. Рядом с соседнем образе поднимается MongoDB для хранения данных пользователя.

3. При удачных регистрации или аутентификации пользователю в ответ отправляется Cookie с token-ом, который в дальнейшем используется для взаимодействия с сервисом.

4. Данные запросы передаются через JSON в Body, а token (jwt) в Cookie.

5. В случае ошибок (неправильный запрос или непредвиденное поведение внутри сервиса) сообщение о том, что пошло не так, отправляется в body ответа, код ошибки проставляется.

6. В базе данных не хранятся пароли в чистом виде, они хешируются с солью.

## Примечания про task_service

1. Используется PostgreSQL в отдельном образе для хранения информации о задачах

2. Используется gRPC. Proto файлы в папке `/task_service/proto`

## Примеры запросов:

### Register

```
curl -v -X POST 'localhost:8080/register' \    
--data '{"username": "name", "password": "kek"}'           
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /register HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Content-Length: 39
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Set-Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w
< Date: Wed, 13 Mar 2024 15:23:08 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

### Authentication

```
curl -v -X POST 'localhost:8080/authenticate' \
--data '{"username": "name", "password": "kek"}'
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /authenticate HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Content-Length: 39
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Set-Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w
< Date: Wed, 13 Mar 2024 15:18:30 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

### Update Profile Information

```
curl -v -X PUT 'localhost:8080/myProfile' \       
--data '{"firstName": "danila", "birthday": "todayepta"}' \
-H 'Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w'
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> PUT /update HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w
> Content-Length: 48
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Wed, 13 Mar 2024 15:36:20 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

### Create Task

```
curl -v -X POST 'localhost:8080/tasks/create' \
--data '{"title": "delete whole arcadia", "description": "yandex should burn in hell", "status": "ready" }' \                                                               -H 'Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w'
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /tasks/create HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w
> Content-Length: 98
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Fri, 05 Apr 2024 21:08:16 GMT
< Content-Length: 15
<
* Connection #0 to host localhost left intact
{"task_id":"1"}
```

### Update Task

```
curl -v -X PUT 'localhost:8080/tasks/update/0' \                                                                               --data '{"title": "dont delete whole arcadia", "description": "yandex should burn in hell", "status": "ready" }' \                                                          -H 'Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w'
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> PUT /tasks/update/052462a1-9143-4679-a54a-e2981d2fba49 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w
> Content-Length: 103
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Fri, 05 Apr 2024 20:28:46 GMT
< Content-Length: 33
<
* Connection #0 to host localhost left intact
Task has been updated succesfully
```

### Delete Task

```
curl -v -X DELETE 'localhost:8080/tasks/delete/0' \
--data '' \                                                                           -H 'Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w'
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> DELETE /tasks/delete/8687c4ed-39ea-4c5c-a881-66dabffc19ba HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w
> Content-Length: 0
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Fri, 05 Apr 2024 20:27:26 GMT
< Content-Length: 33
<
* Connection #0 to host localhost left intact
Task has been deleted succesfully
```

### Get Task

```
curl -v -X GET 'localhost:8080/tasks/get/0' \                                                                                  --data '' \                                                                           -H 'Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w'
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /tasks/get/052462a1-9143-4679-a54a-e2981d2fba49 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w
> Content-Length: 0
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Fri, 05 Apr 2024 20:28:52 GMT
< Content-Length: 97
<
* Connection #0 to host localhost left intact
{"title":"dont delete whole arcadia","description":"yandex should burn in hell","status":"ready"}%
```

### Get TaskList

```
curl -v -X GET 'localhost:8080/tasks/getPage' \
--data '{"offset": 0, "page_size": 2}' \
-H 'Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w'
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /tasks/getPage HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Cookie: token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUifQ.ydrESReHn-a6l9Q2HCXRFlEgx0OA_qB_li1JofY9xCKH3ShpFNuzNDG00J9IVu9Ock3ncLRj0hRCoVHz-sO6gnDZND_xoVqXmw0Kqw1AzeIH7HtHdmFgW--xZIjioHNtp8B2N6VI93kpvz86DBCWo04AhktkiG3rUHcDYdfM-vg0iCopq3EMZh33wmuHIhBUvjqF3NF1ITrofUaJz_R8etwnqpL-diQpY98iKMEoRL9givWsndnYOLex_OKXeGySAJ8SgSDXBvqWlXGFyWOYwnCTweHT-lmsNW6PrWYm1-a83R6WTzVuy31POVMwBSYuiNMdT0Tb2KQprb70NDXv7w
> Content-Length: 29
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Fri, 05 Apr 2024 21:08:58 GMT
< Content-Length: 262
<
* Connection #0 to host localhost left intact
{"tasks":[{"id":{"id":"1"},"task":{"title":"delete whole arcadia","description":"yandex should burn in hell","status":"ready"}},{"id":{"id":"2"},"task":{"title":"delete part of arcadia","description":"yandex should burn in hell","status":"ready"}}],"pageSize":2}
```
