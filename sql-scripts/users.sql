CREATE TABLE IF NOT EXISTS public.users
(
    emailid character varying PRIMARY KEY,
	name character varying COLLATE pg_catalog."default" NOT NULL,
	password character varying COLLATE pg_catalog."default" NOT NULL
)