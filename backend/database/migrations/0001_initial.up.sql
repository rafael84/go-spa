create or replace function public.sync_updated_at ()
returns trigger
as
$body$
begin
    new.updated_at := now ();
    return new;
end;
$body$
language plpgsql volatile cost 100;
