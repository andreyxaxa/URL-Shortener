# URL Shortener - Сервис сокращения URL с аналитикой
Сервис принимает оригинальные адреса и генерирует из них короткие. При переходе на короткий адрес - перенаправляет пользователя на оригинальный URL и собирает статистику по переходам: кто перешёл, когда, с какого устройства и браузера.

[Старт]()

## Обзор

- UI - http://localhost:8080/v1/web
- Документация API - Swagger - http://localhost:8080/swagger
- Конфиг - [config/config.go](https://github.com/andreyxaxa/URL-Shortener/blob/main/config/config.go). Читается из `.env` файла.
- Логгер - [pkg/logger/logger.go](https://github.com/andreyxaxa/URL-Shortener/blob/main/pkg/logger/logger.go). Интерфейс позволяет подменить логгер.
- Кеширование популярных ссылок (Redis) - [internal/repo/cache/link_redis.go](https://github.com/andreyxaxa/URL-Shortener/blob/main/internal/repo/cache/link_redis.go).
- Удобная и гибкая конфигурация HTTP сервера - [pkg/httpserver/options.go](https://github.com/andreyxaxa/URL-Shortener/blob/main/pkg/httpserver/options.go).
  Позволяет конфигурировать сервер в конструкторе таким образом:
  ```go
  httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
  ```
  Аналогичный подход с таким конфигурированием в пакете [pkg/redis](https://github.com/andreyxaxa/URL-Shortener/tree/main/pkg/redis)
- В слое контроллеров применяется версионирование - [internal/controller/restapi/v1](https://github.com/andreyxaxa/URL-Shortener/tree/main/internal/controller/restapi/v1).
  Для версии v2 нужно будет просто добавить папку `restapi/v2` с таким же содержимым, в файле [internal/controller/restapi/router.go](https://github.com/andreyxaxa/URL-Shortener/blob/main/internal/controller/restapi/router.go) добавить строку:
```go
{
		v1.NewLinkRoutes(apiV1Group, lk, l, baseURL+"/v1")
}

{
		v2.NewLinkRoutes(apiV1Group, lk, l, baseURL+"/v2")
}
```
- Graceful shutdown - [internal/app/app.go](https://github.com/andreyxaxa/URL-Shortener/blob/main/internal/app/app.go).

## Запуск

1. Клонируйте репозиторий
2. В корне создайте `.env` файл, скопируйте туда содержимое [env.example](https://github.com/andreyxaxa/URL-Shortener/blob/main/.env.example)
3. Выполните, дождитесь запуска сервиса
   ```
   make compose-up
   ```
4. Перейдите на http://localhost:8080/v1/web и пользуйтесь сервисом.
<img width="1559" height="786" alt="image" src="https://github.com/user-attachments/assets/b27688d7-3854-4be0-8a74-371482685dbf" />


## API

### POST http://localhost:8080/v1/shorten
request:
```json
{
    "url": "https://www.rbc.ru/person/680a00fa9a79477f2728e7a2",
    "custom_alias": "messi"
}
```
response:
```json
{
    "original_url": "https://www.rbc.ru/person/680a00fa9a79477f2728e7a2",
    "short_url": "http://localhost:8080/v1/s/messi"
}
```

### GET http://localhost:8080/v1/s/{short}
request:
```
GET http://localhost:8080/v1/s/messi
```
response:
301 redirect


### GET http://localhost:8080/v1/analytics/{short}
request:
```
GET http://localhost:8080/v1/analytics/messi
```
response:
```json
{
    "analytics": {
        "total_clicks": 12,
        "clicks_by_browser": [
            {
                "browser": "Firefox",
                "clicks": 5
            },
            {
                "browser": "Chrome",
                "clicks": 7
            }
        ],
        "clicks_by_device": [
            {
                "device": "Desktop",
                "clicks": 7
            },
            {
                "device": "Mobile",
                "clicks": 5
            }
        ],
        "recent_clicks": [
            {
                "date": "2026-01-30",
                "clicks": 7
            },
            {
                "date": "2026-01-29",
                "clicks": 5
            }
        ]
    }
}
```

### GET http://localhost:8080/v1/analytics/{short}?group-by=day
request:
```
GET http://localhost:8080/v1/analytics/messi?group-by=day
```
response:
```json
{
    "analytics": {
        "recent_clicks": [
            {
                "date": "2026-01-30",
                "clicks": 7
            },
            {
                "date": "2026-01-29",
                "clicks": 5
            }
        ]
    }
}
```

### GET http://localhost:8080/v1/analytics/{short}?group-by=month
request:
```
GET http://localhost:8080/v1/analytics/messi?group-by=month
```
response:
```json
{
    "analytics": {
        "recent_clicks": [
            {
                "date": "2026-01",
                "clicks": 12
            }
        ]
    }
}
```

### GET http://localhost:8080/v1/analytics/{short}?group-by=browser
request:
```
GET http://localhost:8080/v1/analytics/messi?group-by=browser
```
response:
```json
{
    "analytics": {
        "clicks_by_browser": [
            {
                "browser": "Firefox",
                "clicks": 7
            },
            {
                "browser": "Chrome",
                "clicks": 5
            }
        ]
    }
}
```

### GET http://localhost:8080/v1/analytics/{short}?group-by=device
request:
```
GET http://localhost:8080/v1/analytics/messi?group-by=device
```
response:
```json
{
    "analytics": {
        "clicks_by_device": [
            {
                "device": "Desktop",
                "clicks": 7
            },
            {
                "device": "Mobile",
                "clicks": 5
            }
        ]
    }
}
```

## Прочие `make` команды
Зависимости:
```
make deps
```
docker compose down -v:
```
make compose-down
```
