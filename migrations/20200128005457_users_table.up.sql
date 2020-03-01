CREATE TABLE "users" (
    id bigserial primary key,
    email varchar not null unique,
    enpass varchar not null,
    nickname varchar not null unique,
    first_name varchar not null,
    last_name varchar,    
    avatar varchar not null default '',
    email_verified boolean default false,
    email_verification_code varchar,
    last_login timestamp,
    created timestamp not null default current_timestamp    
);