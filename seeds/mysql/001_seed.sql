DROP TEMPORARY TABLE IF EXISTS seed_numbers;
CREATE TEMPORARY TABLE seed_numbers (n INT PRIMARY KEY);

INSERT INTO seed_numbers (n)
SELECT ones.n + tens.n * 10 + hundreds.n * 100 + 1 AS n
FROM (
    SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
    UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9
) ones
CROSS JOIN (
    SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
    UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9
) tens
CROSS JOIN (
    SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
) hundreds
WHERE ones.n + tens.n * 10 + hundreds.n * 100 + 1 <= 420;

INSERT IGNORE INTO vehicle_class (code, name, description) VALUES
('SEDAN', 'Седан', 'Легковой седан'),
('SUV', 'Кроссовер/SUV', 'Повышенная проходимость'),
('VAN', 'Фургон', 'Грузопассажирский автомобиль'),
('TRUCK', 'Грузовик', 'Грузовой автомобиль'),
('MINIBUS', 'Микроавтобус', 'Пассажирский транспорт'),
('SPECIAL', 'Спецтехника', 'Специализированный транспорт');

INSERT IGNORE INTO vehicle_status (code, name, description) VALUES
('ACTIVE', 'В эксплуатации', 'Используется в работе'),
('MAINTENANCE', 'На обслуживании', 'ТО или ремонт'),
('RENTED_OUT', 'Сдан в аренду', 'Находится у арендатора'),
('RESERVE', 'Резерв', 'Временно не используется'),
('DECOMMISSIONED', 'Списан', 'Выведен из эксплуатации');

INSERT IGNORE INTO fuel_type (code, name, description) VALUES
('AI92', 'Бензин АИ-92', 'Бензин'),
('AI95', 'Бензин АИ-95', 'Бензин'),
('DIESEL', 'Дизель', 'Дизельное топливо'),
('LPG', 'Газ LPG', 'Газовое топливо'),
('ELECTRO', 'Электро', 'Электроэнергия');

INSERT IGNORE INTO transmission_type (code, name, description) VALUES
('MT', 'Механика', 'Механическая КПП'),
('AT', 'Автомат', 'Автоматическая КПП'),
('CVT', 'Вариатор', 'Бесступенчатая КПП'),
('AMT', 'Робот', 'Роботизированная КПП');

INSERT IGNORE INTO acquisition_type (code, name, description) VALUES
('PURCHASE', 'Покупка', 'Покупка в собственность'),
('LEASE', 'Лизинг', 'Лизинговая схема'),
('RENT_IN', 'Аренда в автопарк', 'Аренда стороннего ТС'),
('TRANSFER', 'Передача', 'Внутренняя передача');

INSERT IGNORE INTO maintenance_type (code, name, description) VALUES
('TO', 'Плановое ТО', 'Регламентное обслуживание'),
('REPAIR', 'Ремонт', 'Внеплановый ремонт'),
('DIAGNOSTIC', 'Диагностика', 'Техническая диагностика'),
('BODY', 'Кузовной ремонт', 'Кузовные работы'),
('TIRE', 'Шиномонтаж', 'Работы с шинами');

INSERT IGNORE INTO payment_type (code, name, description) VALUES
('CASH', 'Наличные', 'Оплата наличными'),
('CARD', 'Корпоративная карта', 'Оплата банковской картой'),
('INVOICE', 'Безналичный расчет', 'Оплата по счету');

INSERT IGNORE INTO contract_type (code, name, description) VALUES
('RENT_OUT', 'Аренда (сдача)', 'Сдача ТС в аренду'),
('RENT_IN', 'Аренда (получение)', 'Получение ТС в аренду'),
('SERVICE', 'Сервисный договор', 'Обслуживание и ремонт'),
('PURCHASE', 'Покупка', 'Закупка ТС');

INSERT IGNORE INTO contract_status (code, name, description) VALUES
('DRAFT', 'Черновик', 'Подготовка договора'),
('ACTIVE', 'Действует', 'Активный договор'),
('COMPLETED', 'Завершен', 'Исполненный договор'),
('CANCELLED', 'Отменен', 'Расторгнут');

INSERT IGNORE INTO department (code, name) VALUES
('LOG', 'Логистика'),
('SRV', 'Сервисная служба'),
('SLS', 'Отдел продаж'),
('OPS', 'Операционный отдел'),
('ADM', 'Администрация'),
('REG', 'Региональный филиал');

INSERT IGNORE INTO driver (fio, license_number, phone) VALUES
('Иванов Сергей Петрович', '77АА100001', '+7-901-100-00-01'),
('Петров Андрей Викторович', '77АА100002', '+7-901-100-00-02'),
('Смирнов Николай Олегович', '77АА100003', '+7-901-100-00-03'),
('Кузнецов Роман Игоревич', '77АА100004', '+7-901-100-00-04'),
('Попов Дмитрий Алексеевич', '77АА100005', '+7-901-100-00-05'),
('Соколов Артем Сергеевич', '77АА100006', '+7-901-100-00-06'),
('Лебедев Егор Павлович', '77АА100007', '+7-901-100-00-07'),
('Козлов Максим Юрьевич', '77АА100008', '+7-901-100-00-08'),
('Новиков Михаил Геннадьевич', '77АА100009', '+7-901-100-00-09'),
('Морозов Денис Валерьевич', '77АА100010', '+7-901-100-00-10'),
('Волков Павел Андреевич', '77АА100011', '+7-901-100-00-11'),
('Соловьев Алексей Михайлович', '77АА100012', '+7-901-100-00-12'),
('Васильев Кирилл Николаевич', '77АА100013', '+7-901-100-00-13'),
('Зайцев Илья Артурович', '77АА100014', '+7-901-100-00-14'),
('Павлов Антон Олегович', '77АА100015', '+7-901-100-00-15'),
('Семенов Евгений Викторович', '77АА100016', '+7-901-100-00-16'),
('Голубев Руслан Петрович', '77АА100017', '+7-901-100-00-17'),
('Виноградов Тимур Алексеевич', '77АА100018', '+7-901-100-00-18'),
('Беляев Федор Сергеевич', '77АА100019', '+7-901-100-00-19'),
('Тарасов Константин Ильич', '77АА100020', '+7-901-100-00-20');

INSERT INTO counterparty (type, name, inn, phone, email, address)
SELECT type, name, inn, phone, email, address
FROM (
    SELECT 'company' type, 'ООО ТрансСервис' name, '7701001001' inn, '+7-495-100-10-01' phone, 'info@transservice.ru' email, 'Москва, ул. Первая, д. 10' address
    UNION ALL SELECT 'company', 'ООО СеверЛизинг', '7701001002', '+7-495-100-10-02', 'contact@severleasing.ru', 'Москва, ул. Вторая, д. 20'
    UNION ALL SELECT 'company', 'АО ТехАвто', '7701001003', '+7-495-100-10-03', 'office@techauto.ru', 'Москва, ул. Третья, д. 5'
    UNION ALL SELECT 'company', 'ООО РентПартнер', '7701001004', '+7-495-100-10-04', 'sales@rentpartner.ru', 'Москва, ул. Четвертая, д. 11'
    UNION ALL SELECT 'company', 'ООО АвтоПоставка', '7701001005', '+7-495-100-10-05', 'mail@autopost.ru', 'Москва, ул. Пятая, д. 3'
    UNION ALL SELECT 'company', 'ООО ДорСнаб', '7701001006', '+7-495-100-10-06', 'team@dorsnab.ru', 'Москва, ул. Шестая, д. 7'
    UNION ALL SELECT 'company', 'ООО ГородТранс', '7701001007', '+7-495-100-10-07', 'hello@gortrans.ru', 'Москва, ул. Седьмая, д. 14'
    UNION ALL SELECT 'company', 'ООО СервисМоторс', '7701001008', '+7-495-100-10-08', 'support@servmotors.ru', 'Москва, ул. Восьмая, д. 19'
    UNION ALL SELECT 'individual', 'ИП Захаров И.А.', '7701001009', '+7-926-100-10-09', 'zakharov@ip.ru', 'Московская обл., г. Химки'
    UNION ALL SELECT 'individual', 'ИП Крылов М.С.', '7701001010', '+7-926-100-10-10', 'krylov@ip.ru', 'Московская обл., г. Мытищи'
    UNION ALL SELECT 'individual', 'ИП Романов П.В.', '7701001011', '+7-926-100-10-11', 'romanov@ip.ru', 'Московская обл., г. Балашиха'
    UNION ALL SELECT 'company', 'ООО КонтурАренда', '7701001012', '+7-495-100-10-12', 'rent@contour.ru', 'Москва, ул. Девятая, д. 1'
) t
WHERE NOT EXISTS (SELECT 1 FROM counterparty);

INSERT IGNORE INTO vehicle (
    vin,
    reg_number,
    make,
    model,
    year,
    vehicle_class_id,
    status_id,
    fuel_type_id,
    transmission_type_id,
    acquisition_type_id,
    acquisition_date,
    acquisition_cost,
    current_odometer_km
)
SELECT
    CONCAT('XTA', LPAD(n, 14, '0')) AS vin,
    CONCAT('А', LPAD(n, 3, '0'), 'ВС77') AS reg_number,
    CASE MOD(n, 10)
        WHEN 0 THEN 'Toyota'
        WHEN 1 THEN 'Hyundai'
        WHEN 2 THEN 'Kia'
        WHEN 3 THEN 'Volkswagen'
        WHEN 4 THEN 'Ford'
        WHEN 5 THEN 'Renault'
        WHEN 6 THEN 'LADA'
        WHEN 7 THEN 'Nissan'
        WHEN 8 THEN 'Skoda'
        ELSE 'GAZ'
    END AS make,
    CASE MOD(n, 10)
        WHEN 0 THEN 'Camry'
        WHEN 1 THEN 'Solaris'
        WHEN 2 THEN 'Rio'
        WHEN 3 THEN 'Transporter'
        WHEN 4 THEN 'Transit'
        WHEN 5 THEN 'Duster'
        WHEN 6 THEN 'Largus'
        WHEN 7 THEN 'X-Trail'
        WHEN 8 THEN 'Octavia'
        ELSE 'Gazelle Next'
    END AS model,
    2010 + MOD(n, 14) AS year,
    (SELECT id FROM vehicle_class WHERE code =
        CASE
            WHEN MOD(n, 10) IN (0,1,2,8) THEN 'SEDAN'
            WHEN MOD(n, 10) IN (4,7) THEN 'SUV'
            WHEN MOD(n, 10) IN (3,6,9) THEN 'VAN'
            ELSE 'TRUCK'
        END
    ) AS vehicle_class_id,
    (SELECT id FROM vehicle_status WHERE code =
        CASE
            WHEN n <= 42 THEN 'ACTIVE'
            WHEN n <= 48 THEN 'MAINTENANCE'
            WHEN n <= 54 THEN 'RENTED_OUT'
            WHEN n <= 58 THEN 'RESERVE'
            ELSE 'DECOMMISSIONED'
        END
    ) AS status_id,
    (SELECT id FROM fuel_type WHERE code =
        CASE
            WHEN MOD(n, 10) IN (3,4,9) THEN 'DIESEL'
            WHEN MOD(n, 10) = 7 THEN 'AI92'
            WHEN MOD(n, 10) = 6 THEN 'LPG'
            ELSE 'AI95'
        END
    ) AS fuel_type_id,
    (SELECT id FROM transmission_type WHERE code =
        CASE
            WHEN MOD(n, 4) = 0 THEN 'AT'
            WHEN MOD(n, 4) = 1 THEN 'MT'
            WHEN MOD(n, 4) = 2 THEN 'CVT'
            ELSE 'AMT'
        END
    ) AS transmission_type_id,
    (SELECT id FROM acquisition_type WHERE code =
        CASE
            WHEN MOD(n, 4) = 0 THEN 'PURCHASE'
            WHEN MOD(n, 4) = 1 THEN 'LEASE'
            WHEN MOD(n, 4) = 2 THEN 'RENT_IN'
            ELSE 'TRANSFER'
        END
    ) AS acquisition_type_id,
    DATE_ADD('2018-01-01', INTERVAL MOD(n * 29, 1800) DAY) AS acquisition_date,
    ROUND(850000 + n * 27000, 2) AS acquisition_cost,
    18000 + n * 2300 AS current_odometer_km
FROM seed_numbers
WHERE n <= 60
  AND NOT EXISTS (SELECT 1 FROM vehicle);

INSERT INTO vehicle_assignment (vehicle_id, driver_id, department_id, date_from, date_to, is_primary)
SELECT
    v.id,
    CASE WHEN MOD(v.id, 9) = 0 THEN NULL ELSE MOD(v.id - 1, 20) + 1 END AS driver_id,
    MOD(v.id - 1, 6) + 1 AS department_id,
    DATE_SUB(CURRENT_DATE, INTERVAL (MOD(v.id, 300) + 120) DAY),
    NULL,
    TRUE
FROM vehicle v
WHERE NOT EXISTS (SELECT 1 FROM vehicle_assignment);

INSERT INTO trip_sheet (
    trip_date,
    vehicle_id,
    driver_id,
    department_id,
    odometer_start,
    odometer_end,
    route,
    purpose,
    distance_km
)
SELECT
    DATE_ADD('2025-01-01', INTERVAL MOD(n, 420) DAY) AS trip_date,
    MOD(n - 1, 60) + 1 AS vehicle_id,
    MOD(n - 1, 20) + 1 AS driver_id,
    MOD(n - 1, 6) + 1 AS department_id,
    10000 + n * 17 AS odometer_start,
    10000 + n * 17 + (25 + MOD(n, 140)) AS odometer_end,
    CASE MOD(n, 6)
        WHEN 0 THEN 'Москва - Подольск'
        WHEN 1 THEN 'Москва - Мытищи'
        WHEN 2 THEN 'Москва - Красногорск'
        WHEN 3 THEN 'Москва - Химки'
        WHEN 4 THEN 'Москва - Люберцы'
        ELSE 'Москва - Одинцово'
    END AS route,
    CASE MOD(n, 4)
        WHEN 0 THEN 'Доставка материалов'
        WHEN 1 THEN 'Служебная поездка'
        WHEN 2 THEN 'Выезд к клиенту'
        ELSE 'Перевозка персонала'
    END AS purpose,
    25 + MOD(n, 140) AS distance_km
FROM seed_numbers
WHERE n <= 420
  AND NOT EXISTS (SELECT 1 FROM trip_sheet);

INSERT INTO fuel_txn (txn_ts, vehicle_id, liters, amount, station, odometer_km, payment_type_id)
SELECT
    DATE_ADD('2025-01-01 08:00:00', INTERVAL n HOUR) AS txn_ts,
    MOD(n - 1, 60) + 1 AS vehicle_id,
    ROUND(20 + MOD(n, 45) + MOD(n, 10) * 0.1, 2) AS liters,
    ROUND(1200 + MOD(n, 45) * 95 + MOD(n, 10) * 12, 2) AS amount,
    CASE MOD(n, 5)
        WHEN 0 THEN 'Роснефть'
        WHEN 1 THEN 'Лукойл'
        WHEN 2 THEN 'Газпромнефть'
        WHEN 3 THEN 'Татнефть'
        ELSE 'Shell'
    END AS station,
    11000 + n * 23 AS odometer_km,
    MOD(n - 1, 3) + 1 AS payment_type_id
FROM seed_numbers
WHERE n <= 300
  AND NOT EXISTS (SELECT 1 FROM fuel_txn);

INSERT INTO maintenance_order (
    vehicle_id,
    maintenance_type_id,
    open_date,
    close_date,
    service_name,
    cost,
    description
)
SELECT
    MOD(n - 1, 60) + 1 AS vehicle_id,
    MOD(n - 1, 5) + 1 AS maintenance_type_id,
    DATE_ADD('2025-01-10', INTERVAL (n * 3) DAY),
    DATE_ADD('2025-01-10', INTERVAL (n * 3 + 2 + MOD(n, 7)) DAY),
    CASE MOD(n, 4)
        WHEN 0 THEN 'ООО СервисМоторс'
        WHEN 1 THEN 'АО ТехАвто'
        WHEN 2 THEN 'ООО ДорСнаб'
        ELSE 'ООО ТрансСервис'
    END AS service_name,
    ROUND(3500 + n * 220, 2) AS cost,
    CASE MOD(n, 5)
        WHEN 0 THEN 'Плановая замена масла и фильтров'
        WHEN 1 THEN 'Ремонт ходовой части'
        WHEN 2 THEN 'Диагностика двигателя'
        WHEN 3 THEN 'Замена тормозных колодок'
        ELSE 'Кузовные работы после мелкого ДТП'
    END AS description
FROM seed_numbers
WHERE n <= 120
  AND NOT EXISTS (SELECT 1 FROM maintenance_order);

INSERT IGNORE INTO contract (
    contract_type_id,
    counterparty_id,
    number,
    date_from,
    date_to,
    status_id,
    total_amount,
    notes
)
SELECT
    (SELECT id FROM contract_type WHERE code =
        CASE
            WHEN MOD(n, 4) = 0 THEN 'RENT_OUT'
            WHEN MOD(n, 4) = 1 THEN 'RENT_IN'
            WHEN MOD(n, 4) = 2 THEN 'SERVICE'
            ELSE 'PURCHASE'
        END
    ) AS contract_type_id,
    MOD(n - 1, 12) + 1 AS counterparty_id,
    CONCAT('CNT-2026-', LPAD(n, 4, '0')) AS number,
    DATE_ADD('2025-01-01', INTERVAL (n * 15) DAY),
    DATE_ADD('2025-01-01', INTERVAL (n * 15 + 120 + n * 2) DAY),
    (SELECT id FROM contract_status WHERE code =
        CASE
            WHEN MOD(n, 5) IN (0,1,2) THEN 'ACTIVE'
            WHEN MOD(n, 5) = 3 THEN 'COMPLETED'
            ELSE 'DRAFT'
        END
    ) AS status_id,
    ROUND(150000 + n * 32000, 2) AS total_amount,
    CASE
        WHEN MOD(n, 4) = 0 THEN 'Договор аренды автомобиля'
        WHEN MOD(n, 4) = 1 THEN 'Договор привлечения транспорта'
        WHEN MOD(n, 4) = 2 THEN 'Сервисный договор на обслуживание'
        ELSE 'Закупка транспортного средства'
    END
FROM seed_numbers
WHERE n <= 24
  AND NOT EXISTS (SELECT 1 FROM contract);

INSERT INTO rental_event (
    contract_id,
    vehicle_id,
    pickup_ts,
    return_ts,
    price_per_day,
    deposit,
    notes
)
SELECT
    rc.id,
    vr.id,
    DATE_ADD('2026-01-01 10:00:00', INTERVAL (rc.rn * 12) DAY) AS pickup_ts,
    CASE
        WHEN MOD(rc.rn, 3) = 0 THEN NULL
        ELSE DATE_ADD(DATE_ADD('2026-01-01 10:00:00', INTERVAL (rc.rn * 12) DAY), INTERVAL (7 + MOD(rc.rn, 4)) DAY)
    END AS return_ts,
    ROUND(2600 + rc.rn * 180, 2) AS price_per_day,
    ROUND(12000 + rc.rn * 900, 2) AS deposit,
    'Аренда автомобиля по договору'
FROM (
    SELECT c.id, ROW_NUMBER() OVER (ORDER BY c.id) AS rn
    FROM contract c
    JOIN contract_type ct ON ct.id = c.contract_type_id
    WHERE ct.code = 'RENT_OUT'
) rc
JOIN (
    SELECT v.id, ROW_NUMBER() OVER (ORDER BY v.id) AS rn
    FROM vehicle v
) vr ON vr.rn = rc.rn
WHERE rc.rn <= 12
  AND NOT EXISTS (SELECT 1 FROM rental_event);

UPDATE vehicle v
JOIN (
    SELECT vehicle_id, MAX(odometer_end) AS max_odo
    FROM trip_sheet
    GROUP BY vehicle_id
) t ON t.vehicle_id = v.id
SET v.current_odometer_km = GREATEST(v.current_odometer_km, t.max_odo);

DROP TEMPORARY TABLE IF EXISTS seed_numbers;
