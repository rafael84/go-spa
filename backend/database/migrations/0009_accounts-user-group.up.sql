
create table accounts.user_group
(
    id            bigserial not null primary key,
    state         integer not null default 0,
    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now(),
    user_id       bigint not null,
    group_id      bigint not null,
    json_data     jsonb
)
with ( oids = false );

create index user_group_json_data_idx on accounts.user_group using gin (json_data);

create trigger sync_user_group_updated_at before update 
on accounts.user_group for each row 
execute procedure public.sync_updated_at();

alter table accounts.user_group add foreign key (user_id) references accounts.user(id);
alter table accounts.user_group add foreign key (group_id) references accounts.group(id);
