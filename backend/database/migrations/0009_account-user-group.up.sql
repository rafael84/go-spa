
create table account.user_group
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

create index user_group_json_data_idx on account.user_group using gin (json_data);

create trigger sync_user_group_updated_at before update 
on account.user_group for each row 
execute procedure public.sync_updated_at();

alter table account.user_group add foreign key (user_id) references account.user(id);
alter table account.user_group add foreign key (group_id) references account.group(id);
