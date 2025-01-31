CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(10,2) NOT NULL CHECK (balance >= 0)
);

-- Inserindo os usuários iniciais
INSERT INTO users (id, balance) VALUES (1, 1000.00);
INSERT INTO users (id, balance) VALUES (2, 1000.00);
INSERT INTO users (id, balance) VALUES (3, 1000.00);
INSERT INTO users (id, balance) VALUES (4, 1000.00);
INSERT INTO users (id, balance) VALUES (5, 1000.00);