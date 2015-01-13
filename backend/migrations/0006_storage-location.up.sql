create table storage.location
(
    id           bigserial not null primary key,
    state        integer not null default 0,
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now(),
    name         varchar(100) not null,
    static_url   varchar(240) not null,
    static_path  varchar(240) not null,
    json_data    jsonb
)
with ( oids = false );

create index location_json_data_idx on storage.location using gin (json_data);

create trigger sync_location_updated_at before update 
on storage.location for each row 
execute procedure public.sync_updated_at();
