ALTER TABLE cards 
ADD COLUMN IF NOT EXISTS current_energy INT DEFAULT 5 CHECK (current_energy >= 0 AND current_energy <= 5),
ADD COLUMN IF NOT EXISTS next_refresh_timestamp TIMESTAMP;
