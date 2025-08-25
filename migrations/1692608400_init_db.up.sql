CREATE SEQUENCE IF NOT EXISTS subscription_seq START 1;

CREATE TABLE IF NOT EXISTS subscription (
    id BIGINT PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    monthly_price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);