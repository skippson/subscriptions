# Subscriptions

## Описание
Сервис Subscriptions предназначен для управления записями о подписках пользователей. Для хранения данных используется PostgreSQL

Релизует REST API:

* создание записи о подписке
* получение записи о подписке
* обновление записи о подписке
* удаление записи о подписке
* получение всех записей о подписках
* подсчет суммарной стоимости подписок по месяцам

### Контракт
Реализован Swagger доступный после запуска сервиса по адресу
`${SERVICE_HOST}:${SERVICE_PORT}/api/swagger`

## Локальное развертывание
* Для настройки переменных окружения смотрите `.example.env`

* Запуск
    ```
    docker-compose up --build -d
    ```

* Сервис готов к работе

* Остановка
    ```
    docker-compose down -v
    ```

## Примеры запросов
* Создание записи
    ```
    curl -X POST http://localhost:8080/api/subscriptions \
        -H "Content-Type: application/json" \
        -d '{
                "user_id":"fcb0b8a4-0c3e-4cd0-8659-16a8f6f46f3a",
                "name":"example",
                "price": 200,
                "start_date":"2025-01",
                "end_date":"2026-01"
            }'
    ```

* Получение записи
    ```
    curl -X GET http://localhost:8080/api/subscriptions/e209081e-ab07-485a-86b9-66cef19b24d5
    ```

* Удаление записи
    ```
    curl -X DELETE http://localhost:8080/api/subscriptions/e209081e-ab07-485a-86b9-66cef19b24d5
    ```

* Обновление записи
    ```
    curl -X PATCH http://localhost:8080/api/subscriptions/e209081e-ab07-485a-86b9-66cef19b24d5 \
        -H "Content-Type: application/json" \
        -d '{
                "price": 300,
                "end_date":"2026-02"
            }'
    ```

* Получение списка записей 
    ```
    curl -X GET http://localhost:8080/api/subscriptions/?limit=20&offset=0
    ```

* Подсчет суммарной стоимости подписок
    ```
    curl -X GET http://localhost:8080/api/subscriptions/total_cost?from=2025-05&to=2026-03&user_id=fcb0b8a4-0c3e-4cd0-8659-16a8f6f46f3a
    ```

## Документация
* `config` - Установка конфига
* `internal/adapters/repository/postgres` - Взаимодействия с базой данных
* * `migrations` - Миграции базы данных
* `internal/controllers/http_handlers` - Транспортный слой(реализация запросов)
* * `middleware` - Промежуточная логика
* `internal/models` - Доменные модели
* `internal/server` - Реализация сервера
* `internal/usecase` - Бизнес-логика
* `pkg/logger` - Логгер модель