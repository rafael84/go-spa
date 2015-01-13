create table account.reset_token
(
    id           bigserial not null primary key,
    state        integer not null default 0,
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now(),
    key          varchar(64) not null,
    expiration   timestamptz not null,
    user_id      bigint not null
)
with ( oids = false );

alter table account.reset_token add foreign key(user_id) references account.user(id);

create trigger sync_reset_token_updated_at before update 
on account.reset_token for each row 
execute procedure public.sync_updated_at();
