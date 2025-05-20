# API

0. golang, postman, postgres
1. Прописать api ручки, максимально просто в одном main.go файле, хранить данные в map
Использовать fiber
    1. `POST` `/user` req body{name, age, ...}; resp id(uuid), 201, 400, 500
    2. `GET` `/user/:id`; resp body{id, name, age, ...}, 200, 400, 500
    3. `DELETE` `/user/:id`; resp 200, 400, 500
    4. `PUT` `/user` req body{id, name, age, ...}; resp id(uuid) 200, 400, 500


1. слайс vs массив, структура слайса, что проиходит при append (cap)
2. map, коллизии, эвакуации, swiss table в новой версии go 1.24
3. interface, что под капотом, solid, опп в golang (с примерами чек)
4. Приведение типов
https://www.youtube.com/@Skills_mentor/videos
