-- migrate:up
CREATE TABLE users (
    PRIMARY KEY (id),
    id                    uuid                 DEFAULT gen_random_uuid(),
    auth_id               uuid UNIQUE NOT NULL,
    email                 TEXT UNIQUE NOT NULL,
    first_name            TEXT        NOT NULL,
    last_name             TEXT        NOT NULL,
    phone_number          TEXT UNIQUE,
    email_verified        BOOLEAN              DEFAULT FALSE,
    phone_number_verified BOOLEAN              DEFAULT FALSE,
    created_at            timestamptz NOT NULL DEFAULT NOW(),
    updated_at            timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE INDEX idx_users_id ON users (id);

CREATE INDEX idx_users_auth_id ON users (auth_id);

CREATE TABLE linked_users (
    PRIMARY KEY (id),
    id             uuid                 DEFAULT gen_random_uuid(),
    user_id        uuid REFERENCES users (id),
    linked_user_id uuid REFERENCES users (id),
    created_at     timestamptz NOT NULL DEFAULT NOW(),
    CHECK (user_id <> linked_users.linked_user_id)
);

CREATE INDEX idx_linked_users_id ON linked_users (id);

ALTER TABLE categories
    ADD CONSTRAINT categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id),
    ADD CONSTRAINT categories_linked_users_id_fkey FOREIGN KEY (linked_users_id) REFERENCES linked_users (id);

ALTER TABLE transactions
    ADD CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE accounts
    ADD CONSTRAINT accounts_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id),
    ADD CONSTRAINT accounts_linked_users_id_fkey FOREIGN KEY (linked_users_id) REFERENCES linked_users (id);

DROP TABLE joint_accounts;

-- migrate:down
CREATE TABLE joint_accounts (
    PRIMARY KEY (id),
    id            uuid DEFAULT gen_random_uuid(),
    account_id    uuid NOT NULL REFERENCES accounts (id),
    owner_id      uuid NOT NULL,
    joint_user_id uuid NOT NULL
);

ALTER TABLE accounts
    DROP CONSTRAINT accounts_user_id_fkey,
    DROP CONSTRAINT accounts_linked_users_id_fkey;

ALTER TABLE transactions
    DROP CONSTRAINT transactions_user_id_fkey;

ALTER TABLE categories
    DROP CONSTRAINT categories_user_id_fkey,
    DROP CONSTRAINT categories_linked_users_id_fkey;

DROP TABLE linked_users;

DROP TABLE users;
