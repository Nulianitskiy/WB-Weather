CREATE TABLE city (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    latitude VARCHAR(10),
    longitude VARCHAR(10)
);

CREATE TABLE weather (
    id SERIAL PRIMARY KEY,
    city_id INT NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL,
    temperature NUMERIC(4,2),
    weather_data JSON NOT NULL,
    UNIQUE (city_id, date)
);

INSERT INTO city (name, country, latitude, longitude) VALUES
    ('Москва', 'Россия', '55.7558', '37.6176'),
    ('Санкт-Петербург', 'Россия', '59.9343', '30.3351'),
    ('Новосибирск', 'Россия', '55.0084', '82.9357'),
    ('Екатеринбург', 'Россия', '56.8389', '60.6057'),
    ('Казань', 'Россия', '55.8304', '49.0661'),
    ('Нижний Новгород', 'Россия', '56.2965', '43.9361'),
    ('Челябинск', 'Россия', '55.1644', '61.4368'),
    ('Самара', 'Россия', '53.1959', '50.1000'),
    ('Омск', 'Россия', '54.9893', '73.3682'),
    ('Ростов-на-Дону', 'Россия', '47.2357', '39.7015');
