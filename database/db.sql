create table place (
    id serial primary key ,
    city varchar(100),
    country varchar(100),
    latitude decimal(9,6),
    longitude decimal(9,6)
);

create table weather (
    id serial primary key,
    city_id int not null references place(id) on delete cascade,
    date timestamp,
    temperature decimal(5,3),
    weather_data jsonb not null
)