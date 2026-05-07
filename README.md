# Учёт автопарка

Учебное server-rendered web-приложение на Go для лабораторной работы по базам данных. Система ведёт автомобили, водителей, подразделения, путевые листы, заправки, обслуживание, контрагентов, договоры и аренду.

Приложение работает с двумя СУБД:

- PostgreSQL;
- MySQL.

Активная СУБД выбирается в web-интерфейсе без перезапуска приложения. Все CRUD-страницы, поиски, отчёты и JSON endpoints используют выбранную пользователем базу.

## Стек

- Go 1.24;
- `chi` router;
- `html/template`;
- `database/sql`;
- PostgreSQL через `pgx`;
- MySQL через `github.com/go-sql-driver/mysql`;
- server-rendered HTML, CSS и небольшой JavaScript.

## Быстрый запуск через Docker

```powershell
docker compose up --build
```

Compose поднимает три сервиса:

- `app` на `http://localhost:8080`;
- `postgres` на `localhost:5432`;
- `mysql` на `localhost:3306`.

При старте контейнера приложения выполняются миграции и seed-данные для обеих СУБД.

Данные для входа:

```text
Логин: admin
Пароль: admin123
```

## Локальный запуск без Docker

1. Создайте `.env`:

```powershell
Copy-Item .env.example .env
```

2. Поднимите PostgreSQL и MySQL локально.

3. Создайте базы и пользователя, если они ещё не созданы.

PostgreSQL:

```sql
CREATE USER fleet_user WITH PASSWORD 'fleet_password';
CREATE DATABASE fleet_db OWNER fleet_user;
```

MySQL:

```sql
CREATE DATABASE fleet_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'fleet_user'@'%' IDENTIFIED BY 'fleet_password';
GRANT ALL PRIVILEGES ON fleet_db.* TO 'fleet_user'@'%';
FLUSH PRIVILEGES;
```

4. Примените миграции и seed-данные:

```powershell
go run main.go migrate
go run main.go seed
```

5. Запустите приложение:

```powershell
go run main.go
```

Откройте `http://localhost:8080`.

## Настройки окружения

Основные переменные находятся в `.env.example`.

```env
DEFAULT_DB=postgres

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=fleet_db
POSTGRES_USER=fleet_user
POSTGRES_PASSWORD=fleet_password
POSTGRES_SSLMODE=disable

MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE=fleet_db
MYSQL_USER=fleet_user
MYSQL_PASSWORD=fleet_password
MYSQL_PARSE_TIME=true
MYSQL_LOC=Local
```

Для PostgreSQL сохранена совместимость со старыми `DB_*` переменными. Если указаны `POSTGRES_URL` или `MYSQL_URL`, они имеют приоритет над отдельными параметрами подключения.

## Миграции и seed-данные

PostgreSQL:

- `migrations/001_init.sql`;
- `seeds/001_seed.sql`.

MySQL:

- `migrations/mysql/001_init.sql`;
- `seeds/mysql/001_seed.sql`.

Seed содержит справочники, 60 автомобилей, 20 водителей, подразделения, путевые листы, заправки, ремонты, контрагентов, договоры и события аренды.

## Переключение СУБД

После входа в верхней панели отображается активная база:

```text
Активная БД: PostgreSQL
```

Рядом находится кнопка переключения на вторую СУБД. Выбор хранится в пользовательской session, поэтому разные пользователи могут работать с разными активными БД.

## Что можно демонстрировать

- Полный CRUD основных сущностей.
- Раздел `Справочники` с CRUD для `vehicle_class`, `vehicle_status`, `fuel_type`, `transmission_type`, `acquisition_type`, `maintenance_type`, `payment_type`, `contract_type`, `contract_status`.
- Добавление связанной записи через кнопку `+` рядом с foreign key dropdown.
- Возврат на исходную форму после `+` с восстановлением уже введённых значений и выбором новой записи.
- Поиск автомобилей по марке, модели, госномеру и VIN.
- Каскадный поиск: марка -> модели -> госномера -> VIN.
- Каскадный поиск путевых листов: водитель -> связанные автомобили.
- Отчёты по пробегу, топливу, обслуживанию, статусам, классам, текущей аренде и сводке БД.
- Журнал операций доступен для просмотра и фильтрации.

## JSON endpoints для зависимых списков

- `GET /api/options/vehicles?driver_id=...`
- `GET /api/options/drivers?vehicle_id=...`
- `GET /api/options/models?make=...`
- `GET /api/options/reg-numbers?make=...&model=...`
- `GET /api/options/vins?reg_number=...`
- `GET /api/options/contracts?counterparty_id=...`
- `GET /api/options/rental-events?contract_id=...`
- `GET /api/options/trip-sheets?vehicle_id=...&driver_id=...`
- `GET /api/options/fuel-transactions?vehicle_id=...`
- `GET /api/options/maintenance-orders?vehicle_id=...`

Все endpoints используют активную СУБД из session.

## Проверка

```powershell
go test ./...
docker compose config
```

Если Docker Desktop запущен, полная проверка:

```powershell
docker compose up --build
```

## Компромисс по MySQL

PostgreSQL-миграция использует exclusion constraint для запрета пересечения периодов закрепления автомобиля. В MySQL прямого аналога нет, поэтому такая проверка дополнительно реализована на уровне приложения при создании и редактировании закреплений.
