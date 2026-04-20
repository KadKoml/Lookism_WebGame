ALTER TABLE cards 
DROP COLUMN IF EXISTS current_energy,
DROP COLUMN IF EXISTS next_refresh_timestamp;
