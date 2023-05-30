-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id serial4 NOT NULL,
	"name" varchar(50) NOT NULL,
	nick varchar(50) NOT NULL,
	email varchar(50) NOT NULL,
	"password" varchar(50) NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT users_email_unique UNIQUE (email),
	CONSTRAINT users_nick_unique UNIQUE (nick),
	CONSTRAINT users_pk PRIMARY KEY (id)
);