create table if not exists  clients(
    id serial primary key,
    client_id varchar(100) unique not null,
    client_secret varchar(100) unique not null,
    created_date timestamp not null,
    updated_date timestamp not null,
    isactive boolean default true not null,
    owner integer references users (id)
)