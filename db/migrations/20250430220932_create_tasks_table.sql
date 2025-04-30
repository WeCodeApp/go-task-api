-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clients (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    issuer VARCHAR(255) NOT NULL,
    secret VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL
);

-- Create an index on the deleted field for faster queries
CREATE INDEX idx_tasks_deleted ON tasks(deleted);
CREATE INDEX idx_clients_issuer ON clients(issuer);

-- +goose Down
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS clients;
