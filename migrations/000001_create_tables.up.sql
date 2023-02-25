CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL,
    balance numeric(10, 2) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- Add one user for testing
INSERT INTO users (id, balance, created_at)
VALUES (
        '63e83104-b9a7-4fec-929e-9d08cae3f9b9',
        50,
        now()
    );

CREATE TABLE IF NOT EXISTS transaction (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    amount numeric(10, 2) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    canceled_at timestamp with time zone,
    CONSTRAINT transaction_pkey PRIMARY KEY (id),
    CONSTRAINT transaction_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
