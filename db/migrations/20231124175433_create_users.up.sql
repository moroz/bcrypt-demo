create extension if not exists "citext" with schema "public";

create table users (
  id uuid primary key,
  email citext not null unique,
  password_hash text,
  inserted_at timestamp(0) not null default (now() at time zone 'utc'),
  updated_at timestamp(0) not null default (now() at time zone 'utc')
);
