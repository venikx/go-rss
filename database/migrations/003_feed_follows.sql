-- +goose Up
CREATE TABLE feed_follows (
       id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       user_id INT REFERENCES users(id) ON DELETE CASCADE,
       feed_id INT REFERENCES feeds(id) ON DELETE CASCADE,
       UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
