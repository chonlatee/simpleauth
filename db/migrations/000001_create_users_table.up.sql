create table if not exists users(
    id int(11) not null auto_increment,
    email varchar(100) unique not null,
    username varchar(100) not null,
    password varchar(250) not null,
    created_date  datetime,
    updated_date datetime  default current_timestamp(),
    primary key(id)
)