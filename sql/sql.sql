-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id serial4 NOT NULL,
	"name" varchar(50) NOT NULL,
	nick varchar(50) NOT NULL,
	email varchar(50) NOT NULL,
	"password" varchar(100) NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT users_email_unique UNIQUE (email),
	CONSTRAINT users_nick_unique UNIQUE (nick),
	CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE TABLE public.followers (
	user_id int4 NOT NULL,
	follower_id int4 NOT NULL,
	CONSTRAINT followers_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
	CONSTRAINT followers_follower_id_fkey FOREIGN KEY (follower_id) REFERENCES users(id),
	CONSTRAINT followers_pk PRIMARY KEY (user_id, follower_id)
);