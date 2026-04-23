# Учёт автопарка (Go + PostgreSQL)

Веб-приложение для учебной предметной области **«Учёт автопарка»** (дисциплина «Базы данных»).  
Приложение реализовано как серверный SSR-сайт на Go и PostgreSQL, с обязательными CRUD-операциями, поиском, отчётами, миграциями, seed-данными и запуском через Docker.

## Стек
- Go 1.22+
- PostgreSQL 16+
- `chi` (роутер)
- `html/template` (SSR-шаблоны)
- `database/sql` + `pgx` driver
- Docker + Docker Compose

## Основные возможности
- Авторизация (login/logout, cookie-session, один администратор из env)
- Главная страница с навигацией по всем разделам
- Полный CRUD для сущностей:
  - `vehicle`
  - `driver`
  - `department`
  - `vehicle_assignment`
  - `trip_sheet`
  - `fuel_txn`
  - `maintenance_order`
  - `counterparty`
  - `contract`
  - `rental_event`
- Списки с пагинацией, сортировкой, фильтрами и текстовым поиском
- 3 обязательных поисковых запроса:
  1. Автомобили по марке
  2. Путевые листы по водителю + автомобилю + периоду
  3. Суммарный пробег по водителям и автомобилям за период
- Отчёты:
  - Пробег за период
  - Расходы на топливо за период
  - Расходы на обслуживание/ремонт за период
  - Автомобили по статусам
  - Автомобили по классам
  - Текущие аренды
  - Сводка базы
- Страница справки/о проекте
- Аудит-лог операций INSERT/UPDATE/DELETE (application-level)

## Структура проекта
```text
.
├─ main.go
├─ go.mod
├─ go.sum
├─ .env.example
├─ README.md
├─ Dockerfile
├─ docker-compose.yml
├─ migrations/
├─ seeds/
├─ static/
├─ templates/
└─ internal/
   ├─ config/
   ├─ db/
   ├─ handlers/
   ├─ middleware/
   ├─ models/
   ├─ repositories/
   ├─ services/
   └─ views/
```

## Переменные окружения
Скопируйте пример и при необходимости измените:

```bash
cp .env.example .env
```

Ключевые переменные:
- `APP_HOST`, `APP_PORT`, `APP_ENV`
- `SESSION_KEY`
- `ADMIN_USERNAME`, `ADMIN_PASSWORD`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `DATABASE_URL` (опционально, вместо отдельных DB-параметров)
- `AUTO_MIGRATE`, `AUTO_SEED`

## Локальный запуск (без Docker)
1. Поднимите PostgreSQL и создайте БД (по умолчанию `fleet_db`).
2. Настройте `.env`.
   Используйте параметры БД:
   - `DB_HOST=localhost`
   - `DB_PORT=5432`
   - `DB_USER=postgres`
   - `DB_PASSWORD=Fgrths197+`
   - `DB_NAME=fleet_db`
3. Выполните:

```bash
go mod tidy
go run main.go migrate
go run main.go seed
go run main.go
```

После запуска приложение доступно по адресу: `http://localhost:8080`.

## Команды приложения
- Запуск сервера: `go run main.go` или `go run main.go serve`
- Только миграции: `go run main.go migrate`
- Только seed-данные: `go run main.go seed`

## Docker запуск
```bash
docker compose up --build
```

По умолчанию в `docker-compose.yml`:
- `AUTO_MIGRATE=true`
- `AUTO_SEED=true`

Это означает, что при первом старте схема и демо-данные загрузятся автоматически.

Приложение: `http://localhost:8080`

## Учетные данные по умолчанию
- Логин: `admin`
- Пароль: `admin123`

(Изменяются через `ADMIN_USERNAME` / `ADMIN_PASSWORD`)

## Seed-данные
В `seeds/001_seed.sql` включены реалистичные данные:
- 60 автомобилей (минимум по ТЗ: 50)
- заполненные справочники
- водители, подразделения, контрагенты
- назначения, путевые листы, топливные операции
- обслуживание, договоры, события аренды

## Карта страниц
- `/` — главная
- `/login` — вход
- CRUD-разделы:
  - `/vehicles`
  - `/drivers`
  - `/departments`
  - `/assignments`
  - `/trip-sheets`
  - `/fuel-txns`
  - `/maintenance-orders`
  - `/counterparties`
  - `/contracts`
  - `/rental-events`
- Поиск:
  - `/search`
  - `/search/vehicles-by-make`
  - `/search/trips`
  - `/search/mileage`
- Отчёты:
  - `/reports`
  - `/reports/mileage`
  - `/reports/fuel`
  - `/reports/maintenance`
  - `/reports/by-status`
  - `/reports/by-class`
  - `/reports/current-rentals`
  - `/reports/summary`
- `/help` — справка

## Развертывание на Linux сервере
Рекомендуемый путь:
1. Скопировать репозиторий на сервер.
2. Настроить env (`ADMIN_*`, `SESSION_KEY`, параметры БД).
3. Запустить `docker compose up -d --build`.
4. Проксировать порт приложения (например, через Nginx).
5. Использовать volume `pgdata` для постоянного хранения данных БД.

## Ограничения
- Реализована минимальная модель авторизации: один администратор из env.
- RBAC (admin/viewer) не выделен в отдельную подсистему.


# Учёт автопарка (Go + PostgreSQL)

Веб-приложение для учебной предметной области **«Учёт автопарка»** (дисциплина «Базы данных»).  
Приложение реализовано как серверный SSR-сайт на Go и PostgreSQL, с обязательными CRUD-операциями, поиском, отчётами, миграциями, seed-данными и запуском через Docker.

## Стек
- Go 1.22+
- PostgreSQL 16+
- `chi` (роутер)
- `html/template` (SSR-шаблоны)
- `database/sql` + `pgx` driver
- Docker + Docker Compose

## Основные возможности
- Авторизация (login/logout, cookie-session, один администратор из env)
- Главная страница с навигацией по всем разделам
- Полный CRUD для сущностей:
  - `vehicle`
  - `driver`
  - `department`
  - `vehicle_assignment`
  - `trip_sheet`
  - `fuel_txn`
  - `maintenance_order`
  - `counterparty`
  - `contract`
  - `rental_event`
- Списки с пагинацией, сортировкой, фильтрами и текстовым поиском
- 3 обязательных поисковых запроса:
  1. Автомобили по марке
  2. Путевые листы по водителю + автомобилю + периоду
  3. Суммарный пробег по водителям и автомобилям за период
- Отчёты:
  - Пробег за период
  - Расходы на топливо за период
  - Расходы на обслуживание/ремонт за период
  - Автомобили по статусам
  - Автомобили по классам
  - Текущие аренды
  - Сводка базы
- Страница справки/о проекте
- Аудит-лог операций INSERT/UPDATE/DELETE (application-level)

## Структура проекта
```text
.
├─ main.go
├─ go.mod
├─ go.sum
├─ .env.example
├─ README.md
├─ Dockerfile
├─ docker-compose.yml
├─ migrations/
├─ seeds/
├─ static/
├─ templates/
└─ internal/
   ├─ config/
   ├─ db/
   ├─ handlers/
   ├─ middleware/
   ├─ models/
   ├─ repositories/
   ├─ services/
   └─ views/
```

## Переменные окружения
Скопируйте пример и при необходимости измените:

```bash
cp .env.example .env
```

Ключевые переменные:
- `APP_HOST`, `APP_PORT`, `APP_ENV`
- `SESSION_KEY`
- `ADMIN_USERNAME`, `ADMIN_PASSWORD`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `DATABASE_URL` (опционально, вместо отдельных DB-параметров)
- `AUTO_MIGRATE`, `AUTO_SEED`

## Локальный запуск (без Docker)
1. Поднимите PostgreSQL и создайте БД (по умолчанию `fleet_db`).
2. Настройте `.env`.
   Используйте параметры БД:
   - `DB_HOST=localhost`
   - `DB_PORT=5432`
   - `DB_USER=postgres`
   - `DB_PASSWORD=Fgrths197+`
   - `DB_NAME=fleet_db`
3. Выполните:

```bash
go mod tidy
go run main.go migrate
go run main.go seed
go run main.go
```

После запуска приложение доступно по адресу: `http://localhost:8080`.

## Команды приложения
- Запуск сервера: `go run main.go` или `go run main.go serve`
- Только миграции: `go run main.go migrate`
- Только seed-данные: `go run main.go seed`

## Docker запуск
```bash
docker compose up --build
```

По умолчанию в `docker-compose.yml`:
- `AUTO_MIGRATE=true`
- `AUTO_SEED=true`

Это означает, что при первом старте схема и демо-данные загрузятся автоматически.

Приложение: `http://localhost:8080`

## Учетные данные по умолчанию
- Логин: `admin`
- Пароль: `admin123`

(Изменяются через `ADMIN_USERNAME` / `ADMIN_PASSWORD`)

## Seed-данные
В `seeds/001_seed.sql` включены реалистичные данные:
- 60 автомобилей (минимум по ТЗ: 50)
- заполненные справочники
- водители, подразделения, контрагенты
- назначения, путевые листы, топливные операции
- обслуживание, договоры, события аренды

## Карта страниц
- `/` — главная
- `/login` — вход
- CRUD-разделы:
  - `/vehicles`
  - `/drivers`
  - `/departments`
  - `/assignments`
  - `/trip-sheets`
  - `/fuel-txns`
  - `/maintenance-orders`
  - `/counterparties`
  - `/contracts`
  - `/rental-events`
- Поиск:
  - `/search`
  - `/search/vehicles-by-make`
  - `/search/trips`
  - `/search/mileage`
- Отчёты:
  - `/reports`
  - `/reports/mileage`
  - `/reports/fuel`
  - `/reports/maintenance`
  - `/reports/by-status`
  - `/reports/by-class`
  - `/reports/current-rentals`
  - `/reports/summary`
- `/help` — справка

## Развертывание на Linux сервере
Рекомендуемый путь:
1. Скопировать репозиторий на сервер.
2. Настроить env (`ADMIN_*`, `SESSION_KEY`, параметры БД).
3. Запустить `docker compose up -d --build`.
4. Проксировать порт приложения (например, через Nginx).
5. Использовать volume `pgdata` для постоянного хранения данных БД.

## Ограничения
- Реализована минимальная модель авторизации: один администратор из env.
- RBAC (admin/viewer) не выделен в отдельную подсистему.

