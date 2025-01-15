CREATE TABLE IF NOT EXISTS users (
    id uuid primary key,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    email varchar(255) unique not null,
    phone_number varchar(255) unique not null,
    password varchar(255) not null,
    dob varchar(255),
    gender varchar(255),
    profile_picture varchar(255),
    version integer not null default 0,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now()
)