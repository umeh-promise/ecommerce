CREATE TABLE IF NOT EXISTS products (
    id uuid primary key,
    user_id uuid not null,
    discount integer not null default 0,
    name varchar(255) not null,
    price integer not null,
    description varchar(255) not null,
    image varchar(255), 
    version integer not null default 0,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),

    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

