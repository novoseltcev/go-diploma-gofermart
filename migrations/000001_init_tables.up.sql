BEGIN;

CREATE SEQUENCE IF NOT EXISTS users_seq;
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY DEFAULT nextval('users_seq'::regclass),
    login VARCHAR NOT NULL UNIQUE ,
    hashed_password VARCHAR NOT NULL,
    balance MONEY NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS orders (
    number VARCHAR PRIMARY KEY,
    status VARCHAR NOT NULL DEFAULT 'NEW',
    accrual MONEY NULL,
    user_id INTEGER NOT NULL references users(id),
    uploaded_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NULL
);
CREATE INDEX IF NOT EXISTS ix_orders_user_id ON orders (user_id);

CREATE SEQUENCE IF NOT EXISTS withdrawals_seq;
CREATE TABLE IF NOT EXISTS withdrawals (
    id  INTEGER PRIMARY KEY DEFAULT nextval('withdrawals_seq'::regclass),
    user_id INTEGER NOT NULL references users(id),
    sum MONEY NOT NULL,
    "order" VARCHAR NOT NULL,
    processed_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS ix_withdrawals_user_id ON withdrawals (user_id);

COMMIT;
