create table storage.media
(
    id            bigserial not null primary key,
    state         integer not null default 0,
    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now(),
    name          varchar(140) not null,
    media_type_id bigint not null,
    location_id   bigint not null,
    path          varchar(240) not null,
    json_data     jsonb
)
with ( oids = false );

create index media_json_data_idx on storage.media using gin (json_data);

create trigger sync_media_updated_at before update 
on storage.media for each row 
execute procedure public.sync_updated_at();

alter table storage.media add foreign key (location_id) references storage.location(id);
alter table storage.media add foreign key (media_type_id) references storage.media_type(id);
