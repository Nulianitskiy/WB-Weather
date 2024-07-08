create table place (
    id serial primary key ,
    city varchar(100),
    country varchar(100),
    latitude varchar(10),
    longitude varchar(10)
);

create table weather (
    id serial primary key,
    city_id int not null references place(id) on delete cascade,
    date timestamp,
    temperature numeric(4,2),
    weather_data jsonb not null
)