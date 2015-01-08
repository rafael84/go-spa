create table storage.media_type
(
    id           bigserial not null primary key,
    state        integer not null default 0,
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now(),
    name         varchar(140) not null,
    json_data    jsonb
)
with ( oids = false );

create index media_type_json_data_idx on storage.media_type using gin (json_data);

create trigger sync_media_type_updated_at before update 
on storage.media_type for each row 
execute procedure public.sync_updated_at();
