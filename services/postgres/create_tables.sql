CREATE TABLE IF NOT EXISTS customers (
    id varchar(250) NOT NULL,
    first_name varchar(250) NOT NULL,
    last_name varchar(250) NOT NULL,
    phone_number varchar(250) NOT NULL,
    email varchar(250) NOT NULL,
    created_at BIGINT NOT NULL
);
