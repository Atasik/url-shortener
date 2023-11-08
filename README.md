# url-shortener
## Запуск
1. Настроить env файл:
`````````
POSTGRES_PASSWORD=qwerty
POSTGRES_USER=postgres
POSTGRES_PORT=5432
POSTGRES_HOST=postgres

HTTP_HOST=localhost
GRPC_HOST=localhost
`````````
2. Настроить конфиги под себя (папка configs)
3. В зависимости от того, что запускается, выполнить команды:
- Для запуска inmemory решения:
`````````
make run
`````````
- Для запуска Postgresql решения:
`````````
make run-db
`````````
А также затем применить миграции:
`````````
make migrate
`````````
4. По-идее всё должно работать. Если оставили 8080 порт, то ещё можно будет swagger потыкать http://localhost:8080/swagger/

## Curl-запросы
1. Создание токена
`````````
curl -X 'POST' \
  'http://localhost:8080/api/v1/link' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "original_url": "https://vk.com/dtf"
`````````

2. Получение оригинальной ссылки
`````````
curl -X 'GET' \
  'http://localhost:8080/api/v1/link/aaaaaaaaab' \
  -H 'accept: application/json'
`````````

## Формат ответов:
- Оригинальная ссылка:
`````````
{
  "link": "https://vk.com/dtf"
}
`````````
- Токен:
`````````
{
  "token": "blabla1234"
}
`````````

## Ещё пару слов
Сам алгоритм, конечно, вообще не безопасен. Но я именно отталкивался от задания, где про безопасность ничего не сказано, но зато сказано про отсутствие условных коллизий, что я и реализовал.
Так бы я условно взял md5 биты и их конвертировал их, но тогда бы с каждой новой ссылкой вероятность получить коллизию бы возрастала.
