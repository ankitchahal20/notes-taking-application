CREATE TABLE IF NOT EXISTS public.notes
(
    id serial PRIMARY KEY,
    note character varying COLLATE pg_catalog."default" NOT NULL
)