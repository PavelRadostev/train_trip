# Microservice train_trip

**АХТУНГ!!!** Репозиторий содержит ВСЕ файлы, включая чувственные данные (тестовые), чтобы запуститься на удаленной машине.

### maintainer

| Команда  |   Ответственный    |                        Email |
|----------|:------------------:|-----------------------------:|
| Бизнес 1 |  Павел Радостев    |                              |

___

Сервис хранения событий погрузки/разгрузки состава

Входящие события:


Исходящие события: нет.

Используемые запросы: нет.

Реализуемые запросы:

- TrainLoadQuery (погрузки)
- TrainUnloadQuery (разгрузки)

Реализуемые команды:

## Installation
### Apply all pending migrations
`go run cmd/migrate/main.go up`

Rollback last migration
`go run cmd/migrate/main.go down`

Check current version
`go run cmd/migrate/main.go version`

## Installation


### Dependencies


### Linux dependencies


### Локальный запуск compose.


### Локальная разработка.


## Documentation



## Git
git branch -M master
git add -A
git status
git commit -m "Init MS train_trip"
git push origin master

## Tests


## Сборка docker контейнера


## Запуск docker контейнера


## Подключение к запущенному docker контейнеру

