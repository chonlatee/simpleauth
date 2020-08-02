create table if not exists users(
    id serial primary key,
    email varchar(100) unique not null,
    password varchar(100) not null,
    isactive boolean default true not null,
    created_date timestamp not null,
    updated_date timestamp not null
)