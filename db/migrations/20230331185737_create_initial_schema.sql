-- migrate:up
CREATE FUNCTION update_updated_at() RETURNS TRIGGER AS
$update_updated_at$
BEGIN
    new.updated_at = NOW();
    RETURN new;
END;
$update_updated_at$ LANGUAGE plpgsql;

CREATE TYPE cycle_type AS ENUM (
    'weekly',
    'monthly',
    'quarterly',
    'semiannually',
    'annually'
    );

CREATE TYPE transaction_type AS ENUM (
    'charge',
    'refund',
    'deposit',
    'withdrawal',
    'interest',
    'adjustment'
    );

CREATE TYPE account_type AS ENUM (
    'checking',
    'savings',
    'credit'
    );

CREATE TYPE currency AS ENUM (
    'USD'
    );

CREATE TABLE accounts (
    PRIMARY KEY (id),
    id                 uuid                  DEFAULT gen_random_uuid(),
    user_id            uuid         NOT NULL,
    account_type       account_type NOT NULL,
    name               TEXT         NOT NULL,
    account_number     BIGINT,
    routing_number     BIGINT,
    initial_balance    BIGINT,
    credit_limit       BIGINT,
    statement_end_date DATE,
    payment_due_date   DATE,
    linked_users_id    uuid,
    created_at         timestamptz  NOT NULL DEFAULT NOW(),
    updated_at         timestamptz  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts (user_id);

CREATE TRIGGER update_accounts_updated_at
    BEFORE UPDATE
    ON accounts
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();


CREATE TABLE joint_accounts (
    PRIMARY KEY (id),
    id            uuid DEFAULT gen_random_uuid(),
    account_id    uuid NOT NULL REFERENCES accounts (id),
    owner_id      uuid NOT NULL,
    joint_user_id uuid NOT NULL
);

CREATE TABLE categories (
    PRIMARY KEY (id),
    id                 uuid                 DEFAULT gen_random_uuid(),
    user_id            uuid        NOT NULL,
    name               TEXT        NOT NULL,
    parent_category_id uuid REFERENCES categories (id),
    allowance          BIGINT,
    cycle_type         cycle_type,
    rollover           BOOLEAN              DEFAULT FALSE,
    linked_users_id    uuid,
    created_at         timestamptz NOT NULL DEFAULT NOW(),
    updated_at         timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_user_id ON categories (user_id);
CREATE INDEX idx_categories_linked_users_id ON categories (linked_users_id);
CREATE INDEX idx_categories_parent_category_id ON categories (parent_category_id);

CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE
    ON categories
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TABLE transactions (
    id               uuid PRIMARY KEY          DEFAULT gen_random_uuid(),
    user_id          uuid             NOT NULL,
    category_id      uuid REFERENCES categories (id),
    transaction_date timestamptz      NOT NULL DEFAULT NOW(),
    description      TEXT,
    transaction_type transaction_type NOT NULL,
    account_id       uuid REFERENCES accounts (id),
    amount           BIGINT,
    currency         currency         NOT NULL DEFAULT 'USD',
    comment          TEXT,
    created_at       timestamptz      NOT NULL DEFAULT NOW(),
    updated_at       timestamptz      NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_category_id ON transactions (category_id);
CREATE INDEX idx_transactions_account_id ON transactions (account_id);
CREATE INDEX idx_transactions_user_id ON transactions (user_id);

CREATE TRIGGER update_transaction_updated_at
    BEFORE UPDATE
    ON transactions
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- migrate:down
DROP TABLE transactions;
DROP TABLE categories;
DROP TABLE accounts;
DROP TABLE joint_accounts;

DROP FUNCTION update_updated_at();