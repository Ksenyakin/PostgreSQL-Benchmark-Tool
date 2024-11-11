CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active'
    );

INSERT INTO users (name, email, status)
SELECT
    'User' || generate_series(1, 1000) AS name,
    'user' || generate_series(1, 1000) || '@example.com' AS email,
    (ARRAY['active', 'inactive', 'suspended'])[floor(random() * 3 + 1)] AS status;
