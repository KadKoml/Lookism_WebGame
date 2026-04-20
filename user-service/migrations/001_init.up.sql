CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cards (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id VARCHAR(10) NOT NULL,
    level INT DEFAULT 1 CHECK (level >= 1 AND level <= 60),
    merge_stars INT DEFAULT 0 CHECK (merge_stars >= 0)
);

CREATE TABLE IF NOT EXISTS squads (
    user_id VARCHAR(36) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    card_id_1 VARCHAR(36) REFERENCES cards(id) ON DELETE SET NULL,
    card_id_2 VARCHAR(36) REFERENCES cards(id) ON DELETE SET NULL,
    card_id_3 VARCHAR(36) REFERENCES cards(id) ON DELETE SET NULL
);

CREATE INDEX idx_cards_user_id ON cards(user_id);
