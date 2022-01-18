CREATE TABLE IF NOT EXISTS commands(
	id serial PRIMARY KEY,
	issuer text NOT NULL,
	argv text[],
	description text,
	status integer,
	std_out bytea,
	std_err bytea,
	create_time timestamp with time zone,
	update_time timestamp with time zone,
	delete_time timestamp with time zone,
	start_time timestamp with time zone,
	end_time timestamp with time zone
);
