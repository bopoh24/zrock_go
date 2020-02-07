CREATE TABLE "user" (
    id bigserial primary key,
    email varchar not null unique,
    enpass varchar not null,
    nick varchar not null unique,
    first_name varchar not null,
    last_name varchar,    
    avatar varchar,
    last_login timestamp,
    created timestamp not null default current_timestamp    
);