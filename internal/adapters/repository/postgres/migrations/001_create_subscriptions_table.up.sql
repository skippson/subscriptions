create extension if not exists "pgcrypto";

create table if not exists subscriptions (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    name text not null,
    price int not null,
    start_date date not null,
    end_date date,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create index idx_subscriptions_user_id_name on subscriptions(user_id,name);