# PROJECT REQUIREMENTS FOR CODEX

## 0. Status of this document
This document is the binding technical specification for the project. It is based on:
1. the uploaded lab report on the selected subject area: fleet management / vehicle fleet accounting;
2. the uploaded lab manual for the database course, especially the requirements of lab work №3;
3. the explicit project constraints chosen by the student: use PostgreSQL instead of MySQL, use Go instead of PHP/JSP, provide a web UI, make the project deployable on a server, and support Docker and terminal launch.

If the original lab manual and this specification conflict in technology choice, this specification wins for the technology stack, while the functional/educational requirements of the lab manual remain mandatory.

---

## 1. Project goal
Implement a complete web application for the course "Databases" for the subject area:

**"Fleet management / vehicle fleet accounting"**

The application must provide web access to a relational database and support:
- creation, viewing, editing, and deletion of data;
- search queries;
- statistical / aggregate reports;
- a help/about page for the database;
- local launch from terminal through `main.go`;
- Docker-based launch;
- readiness for deployment to a Linux server.

This project is an adaptation of lab work №3. The original manual demonstrates PHP+MySQL or JSP+Oracle, but this project must be implemented with:
- **PostgreSQL** as DBMS;
- **Go** as backend language;
- **server-rendered web interface** in browser.

---

## 2. Hard source requirements that MUST be preserved
These are binding functional requirements inherited from the uploaded lab materials.

### 2.1. Web access to DB is mandatory
The result must be a **web application** that works through a browser and gives access to the database.

### 2.2. The application must repeat the educational logic of the lab
The web application must preserve the functional meaning of the course assignment:
- data input;
- data search;
- statistical / summary information about the database.

### 2.3. One main page is mandatory
There must be **one main page** that contains navigation/buttons/links to:
- input/edit forms for all major entities;
- search pages;
- report/statistics pages;
- help/about page;
- exit/logout if authentication is implemented.

Every auxiliary page should allow returning back to the main page or relevant list page.

### 2.4. Search must cover three educational categories
At minimum, the project must implement the following three query classes:
1. search by one attribute from one table;
2. search by synchronized attributes from two or more related tables;
3. aggregate/statistical query over the DB.

### 2.5. Database scale requirements
The project must contain:
- at least **5 tables**;
- at least **25 attributes total** across the implemented schema;
- at least **50 records in the main table**.

For this project, the main table is `vehicle`.

### 2.6. A help/about page is mandatory
A separate page with concise help/about information about the DB and the project is required.

### 2.7. Demonstration readiness is mandatory
The final project must be usable for demonstration to the teacher in a browser, with real dynamic interaction with the DB.

---

## 3. Binding adaptation for this project
These are not optional. They are project-specific overrides chosen for implementation.

### 3.1. Technology stack
Mandatory stack:
- **Go 1.22+**
- **PostgreSQL 16+**
- HTML templates rendered on server side
- CSS for basic styling
- optional small vanilla JS only where needed
- Docker / Docker Compose

### 3.2. Forbidden substitutions
Do NOT use:
- MySQL;
- MariaDB;
- PHP;
- JSP;
- Oracle;
- Access;
- heavy frontend SPA frameworks;
- heavy ORM-first architecture.

### 3.3. DB access style
Use clear SQL access. Preferred options:
- `pgx` or `database/sql` with PostgreSQL driver;
- explicit SQL queries in repository layer.

Avoid hiding the relational logic behind a complex ORM.

### 3.4. Application style
The project must be a classical web app with server-side rendering. The browser UI must be sufficient to manage the database without external admin tools.

---

## 4. Recommended implementation choices that become binding for this repo
To avoid ambiguity, implement exactly this unless there is a hard technical blocker.

### 4.1. HTTP/router
Use:
- `chi` router for HTTP routing.

### 4.2. Templates
Use:
- standard `html/template`.

### 4.3. Configuration
Use:
- environment variables;
- optional `.env` support;
- `.env.example` in repository root.

### 4.4. Migrations
Use SQL migrations located in `migrations/`.
If a migration tool is used, keep migrations plain SQL and document execution clearly.

### 4.5. Seed data
Provide seed data scripts/files to populate the DB with realistic demo data.

### 4.6. Logging
Provide application logging for startup, shutdown, DB errors, and HTTP/internal failures.

---

## 5. Subject area that MUST be implemented
Subject area: **vehicle fleet accounting**.

The system must model the following real-world processes:
- vehicle cards;
- driver cards;
- departments;
- assignment of vehicle to driver/department over time;
- trip sheets / mileage tracking;
- fuel transactions;
- maintenance and repair orders;
- counterparties;
- contracts;
- rental issue/return events;
- system dictionaries/reference tables.

The system must support purchased vehicles, leased vehicles, rented-in vehicles, and rent-out scenarios.

---

## 6. Mandatory database schema
Implement at least the following tables.

### 6.1. Reference tables
Must exist as separate tables:
- `vehicle_class`
- `vehicle_status`
- `fuel_type`
- `transmission_type`
- `acquisition_type`
- `maintenance_type`
- `payment_type`
- `contract_type`
- `contract_status`

Each reference table should contain at least:
- `id`
- `code`
- `name`
- optional `description`

### 6.2. Main business tables
#### `vehicle`
Fields:
- `id`
- `vin`
- `reg_number`
- `make`
- `model`
- `year`
- `vehicle_class_id`
- `status_id`
- `fuel_type_id`
- `transmission_type_id`
- `acquisition_type_id`
- `acquisition_date`
- `acquisition_cost`
- `current_odometer_km`
- `created_at`
- `updated_at`

#### `driver`
Fields:
- `id`
- `fio`
- `license_number`
- `phone`
- `created_at`
- `updated_at`

#### `department`
Fields:
- `id`
- `code`
- `name`
- `created_at`
- `updated_at`

#### `vehicle_assignment`
Fields:
- `id`
- `vehicle_id`
- `driver_id` nullable
- `department_id`
- `date_from`
- `date_to` nullable
- `is_primary`
- `created_at`

#### `trip_sheet`
Fields:
- `id`
- `trip_date`
- `vehicle_id`
- `driver_id`
- `department_id`
- `odometer_start`
- `odometer_end`
- `route`
- `purpose`
- `distance_km`
- `created_at`

#### `fuel_txn`
Fields:
- `id`
- `txn_ts`
- `vehicle_id`
- `liters`
- `amount`
- `station`
- `odometer_km`
- `payment_type_id`
- `created_at`

#### `maintenance_order`
Fields:
- `id`
- `vehicle_id`
- `maintenance_type_id`
- `open_date`
- `close_date` nullable
- `service_name`
- `cost`
- `description`
- `created_at`

#### `counterparty`
Fields:
- `id`
- `type`
- `name`
- `inn` nullable
- `phone` nullable
- `email` nullable
- `address` nullable
- `created_at`
- `updated_at`

#### `contract`
Fields:
- `id`
- `contract_type_id`
- `counterparty_id`
- `number`
- `date_from`
- `date_to`
- `status_id`
- `total_amount` nullable
- `notes` nullable
- `created_at`
- `updated_at`

#### `rental_event`
Fields:
- `id`
- `contract_id`
- `vehicle_id`
- `pickup_ts`
- `return_ts` nullable
- `price_per_day`
- `deposit` nullable
- `notes` nullable
- `created_at`
- `updated_at`

### 6.3. Audit log
Implement table:
- `audit_log`

Fields:
- `id`
- `ts`
- `actor`
- `action`
- `entity`
- `entity_id`
- `details_before` jsonb nullable
- `details_after` jsonb nullable

Minimum expectation:
- log INSERT/UPDATE/DELETE for key business tables;
- if full trigger coverage is too large, at least provide application-level logging for all CRUD operations and document it.

---

## 7. Mandatory integrity constraints
These constraints are binding.

### 7.1. vehicle
- `vin` is required, unique, length 17;
- `reg_number` is required, unique;
- `year` must be within a reasonable range;
- `acquisition_cost >= 0`;
- `current_odometer_km >= 0`.

### 7.2. driver
- `license_number` is required and unique.

### 7.3. vehicle_assignment
- `date_to >= date_from` when `date_to` is not null;
- overlapping assignment periods for the same vehicle must be prevented.

### 7.4. trip_sheet
- `odometer_end >= odometer_start`;
- `distance_km >= 0`;
- `distance_km` should be derived from odometer values or validated against them.

### 7.5. fuel_txn
- `liters > 0`;
- `amount >= 0`;
- `odometer_km >= 0`.

### 7.6. maintenance_order
- `close_date >= open_date` when `close_date` is not null;
- `cost >= 0`.

### 7.7. contract
- `number` should be unique;
- `date_to >= date_from`;
- `total_amount >= 0` when not null.

### 7.8. rental_event
- `return_ts >= pickup_ts` when `return_ts` is not null;
- `price_per_day >= 0`;
- `deposit >= 0` when not null.

### 7.9. Foreign keys
All relations must be implemented with proper foreign keys.

---

## 8. Mandatory relations
Implement at least these relations:
- `vehicle -> trip_sheet` (1:N)
- `driver -> trip_sheet` (1:N)
- `department -> trip_sheet` (1:N)
- `vehicle -> fuel_txn` (1:N)
- `vehicle -> maintenance_order` (1:N)
- `vehicle -> vehicle_assignment` (1:N)
- `counterparty -> contract` (1:N)
- `contract -> rental_event` (1:N)
- `vehicle -> rental_event` (1:N)

---

## 9. Mandatory UI pages
Implement these browser pages.

### 9.1. Main page / dashboard
Must contain navigation to:
- vehicles
- drivers
- departments
- assignments
- trip sheets
- fuel transactions
- maintenance orders
- counterparties
- contracts
- rental events
- search
- reports/statistics
- help/about
- login/logout if authentication is enabled

### 9.2. CRUD pages
For each major business table provide:
- list page;
- create page;
- detail page or card page;
- edit page;
- delete action with confirmation.

At minimum full CRUD must exist for:
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

### 9.3. Form behavior
Forms must:
- validate required fields;
- display readable validation errors;
- preserve user input on validation failure when practical;
- use dropdowns/selects for foreign keys;
- provide navigation back to list or home.

### 9.4. Lists
List pages must support:
- pagination;
- sorting;
- basic filtering;
- text search where relevant.

---

## 10. Mandatory search pages and queries
These three searches are mandatory because they directly match the uploaded subject-area report.

### 10.1. Search 1: one-table thematic search
**Search vehicles by make**.

Input:
- vehicle make

Output columns:
- make
- model
- year
- reg_number
- status_name

### 10.2. Search 2: synchronized search across related tables
**Trip sheets by driver and vehicle for a period**.

Input:
- driver
- vehicle registration number
- start date
- end date

Output columns:
- trip_date
- driver fio
- reg_number
- odometer_start
- odometer_end
- distance_km

### 10.3. Search 3: aggregate query
**Total mileage by drivers and vehicles for a period**.

Input:
- start date
- end date

Output columns:
- driver fio
- vehicle make
- vehicle model
- total_km

These pages must be visible in the UI and runnable without admin tools.

---

## 11. Mandatory reports
Implement at least the following report/statistics pages:
- mileage for a period;
- fuel expenses for a period;
- maintenance and repair expenses for a period;
- vehicles by status;
- vehicles by class;
- vehicles currently in rental;
- database summary/help page.

The report pages may be table-based. Charts are optional.

---

## 12. Database help/about page
Implement a dedicated page that includes:
- project title;
- subject area description;
- what entities exist;
- what the main table is;
- what searches and reports are implemented;
- short usage instructions.

This page should be written in clear Russian.

---

## 13. Seed data requirements
Provide realistic demo data.

Mandatory minimums:
- `vehicle`: at least **50** rows;
- all reference tables populated;
- enough rows in other tables to demonstrate every search/report page.

Seed data must be coherent, not random nonsense.
For example:
- realistic VIN-like values;
- realistic registration numbers;
- plausible makes/models/years;
- assignments, trips, fuel operations, maintenance events, rental events linked consistently.

---

## 14. Authentication and authorization
Since the project is intended for server deployment, implement at least minimal authentication.

Minimum acceptable implementation:
- login page;
- one admin user from environment variables or seed;
- session/cookie based auth;
- write operations protected from anonymous access.

Preferred:
- role `admin`;
- optional role `viewer` with read-only access.

If a full role model is too expensive, implement one admin account and document the limitation.

---

## 15. Project structure
Use a clean structure similar to this:

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

Exact names may differ slightly, but the architecture must remain modular and understandable.

---

## 16. Launch requirements
Two launch modes are mandatory.

### 16.1. Local launch without Docker
The project must launch from terminal through:
- `go run main.go`

or an equally simple command, but `main.go` must remain the visible entry point.

### 16.2. Docker launch
Provide:
- `Dockerfile` for application;
- `docker-compose.yml` with app + PostgreSQL;
- persistent volume for PostgreSQL data;
- environment variable wiring;
- documented migration/seed execution.

Preferred:
- app container waits for DB readiness;
- healthcheck;
- production-friendly defaults.

---

## 17. Server deployment requirements
The project must be deployable to a Linux server.

At minimum:
- document production launch via Docker Compose;
- support configuration via environment variables;
- expose application port;
- be compatible with reverse proxy setup;
- store DB data in volume;
- provide deployment instructions in README.

---

## 18. README requirements
`README.md` must include:
- project overview;
- stack;
- local setup;
- Docker setup;
- environment variables;
- migration steps;
- seed steps;
- default login credentials or how to configure them;
- page map / available features;
- deployment notes.

---

## 19. UX requirements
The UI does not need to be fancy, but it must be clean and usable.

Requirements:
- Russian language UI;
- readable tables and forms;
- consistent navigation;
- validation/error messages in Russian;
- no broken links;
- no placeholder-only pages.

---

## 20. Quality requirements
Mandatory quality bar:
- the project must compile;
- routes must work;
- templates must render;
- migrations must apply cleanly on empty DB;
- seed data must load successfully;
- CRUD operations must work at least for the main entities;
- required search/report pages must execute against real DB data.

Add at least minimal tests where practical:
- unit tests for validation or service logic;
- repository tests are optional;
- integration tests optional but welcome.

---

## 21. What is NOT optional
The following are mandatory and cannot be skipped:
- PostgreSQL
- Go backend
- browser-based web UI
- SQL migrations
- seed/demo data
- main page
- CRUD for main entities
- three mandatory search scenarios
- report/statistics pages
- help/about page
- launch via `main.go`
- Docker launch
- README

---

## 22. What may be simplified if necessary
If time is limited, these may be simplified but not fully removed:
- audit log may be partial rather than full trigger-based;
- role model may be minimal;
- styling may be basic;
- charts may be omitted if table reports exist.

However, do not simplify away the core educational requirements.

---

## 23. Final delivery checklist
The final repository must contain:
- source code of the Go application;
- SQL migrations;
- seed files/scripts;
- HTML templates;
- static assets;
- Dockerfile;
- docker-compose.yml;
- `.env.example`;
- README;
- all code needed to run the project end-to-end.

The project is considered complete only if a reviewer can:
1. create/start PostgreSQL;
2. apply migrations;
3. load seed data;
4. run the app;
5. open it in browser;
6. log in;
7. create/edit/delete records;
8. run the required searches;
9. open reports;
10. confirm at least 50 vehicles exist.

