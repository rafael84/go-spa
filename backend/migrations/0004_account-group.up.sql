create table account.group
(
    id           bigserial not null primary key,
    state        integer not null default 0,
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now(),
    name         varchar(140) not null,
    json_data    jsonb
)
with ( oids = false );

create index group_json_data_idx on account.group using gin (json_data);

create trigger sync_group_updated_at before update 
on account.group for each row 
execute procedure public.sync_updated_at();
