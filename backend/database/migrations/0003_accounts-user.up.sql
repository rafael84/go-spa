create table accounts.user
(
    id           bigserial not null primary key,
    state        integer not null default 0,
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now(),
    email        varchar(140) not null,
    password     varchar(30) not null,
    json_data    jsonb
)
with ( oids = false );

create index user_json_data_idx on accounts.user using gin (json_data);

create trigger sync_user_updated_at before update 
on accounts.user for each row 
execute procedure public.sync_updated_at();
