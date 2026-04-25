CREATE TABLE IF NOT EXISTS currency (
    user_id VARCHAR(36) PRIMARY KEY,
    balance INT DEFAULT 1000 CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS energy (
    user_id VARCHAR(36) PRIMARY KEY,
    current_energy INT DEFAULT 40 CHECK (current_energy >= 0),
    max_energy INT DEFAULT 40,
    last_refill TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS daily_rewards (
    user_id VARCHAR(36),
    day INT CHECK (day >= 1 AND day <= 7),
    description VARCHAR(100),
    claimed BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (user_id, day)
);

CREATE TABLE IF NOT EXISTS quests (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    description VARCHAR(200),
    is_completed BOOLEAN DEFAULT FALSE,
    is_claimed BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_quests_user_id ON quests(user_id);
