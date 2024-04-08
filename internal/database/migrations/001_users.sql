-- +goose Up
CREATE TABLE users (
       id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
       name VARCHAR(255) NOT NULL,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
               encode(sha256(random()::text::bytea), 'hex')
       )
);

-- +goose Down
DROP TABLE users;
