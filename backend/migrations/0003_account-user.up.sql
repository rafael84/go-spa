create table account.user
(
    id           bigserial not null primary key,
    state        integer not null default 0,
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now(),
    email        varchar(140) not null,
    password     varchar(240) not null,
    json_data    jsonb
)
with ( oids = false );

create index user_json_data_idx on account.user using gin (json_data);

create trigger sync_user_updated_at before update 
on account.user for each row 
execute procedure public.sync_updated_at();
