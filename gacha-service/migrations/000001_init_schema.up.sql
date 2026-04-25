CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL DEFAULT 'test@example.com',
    last_seen_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    active_squad_ids UUID[] DEFAULT '{}'
);

CREATE TABLE user_wallets (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    money_balance INT NOT NULL DEFAULT 0,
    energy_balance INT NOT NULL DEFAULT 100,
    max_energy INT NOT NULL DEFAULT 100,
    energy_updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE card_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    faction VARCHAR(255) NOT NULL,
    base_hp INT NOT NULL,
    base_str INT NOT NULL,
    base_mana INT NOT NULL,
    base_agility INT NOT NULL,
    base_reaction INT NOT NULL,
    base_durability INT NOT NULL,
    base_power DECIMAL(4,2) NOT NULL,
    base_speed DECIMAL(4,2) NOT NULL,
    base_technique DECIMAL(4,2) NOT NULL,
    passive_skill JSONB NOT NULL,
    active_skill JSONB NOT NULL,
    is_boss BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE user_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    template_id UUID REFERENCES card_templates(id) ON DELETE CASCADE,
    level INT NOT NULL DEFAULT 1,
    merge_stars INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (user_id, template_id)
);

CREATE TABLE shop_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    cost INT NOT NULL,
    description TEXT
);

-- Seed Data: 16 Lookism Characters
INSERT INTO card_templates (name, faction, base_hp, base_str, base_mana, base_agility, base_reaction, base_durability, base_power, base_speed, base_technique, passive_skill, active_skill, is_boss) VALUES
('Daniel Park', 'J High', 430, 210, 4, 48, 50, 42, 1.35, 1.35, 1.45, '{"name": "Adaptation", "effect": "tech_buff", "value": 0.10, "trigger": "round_1"}', '{"name": "Adapt Strike", "cost": 2, "damage_multiplier": 1.20, "buff": {"stat": "technique", "value": 0.10}}', FALSE),
('Zack Lee', 'J High', 390, 225, 2, 42, 44, 50, 1.45, 1.20, 1.20, '{"name": "Boxing Combo", "effect": "damage_buff_chance", "chance": 0.20, "value": 0.25}', '{"name": "Boxing Rush", "cost": 1, "hits": 2, "damage_multiplier": 0.65}', FALSE),
('Vasco', 'J High', 470, 240, 2, 30, 35, 58, 1.55, 1.10, 1.15, '{"name": "Hero Strike", "effect": "str_buff_low_hp", "hp_threshold": 0.50, "value": 0.20}', '{"name": "Hero Punch", "cost": 2, "damage_multiplier": 1.50}', FALSE),
('Jay Hong', 'J High', 360, 190, 5, 55, 48, 35, 1.20, 1.50, 1.40, '{"name": "Weapon Mastery", "effect": "squad_tech_buff", "value": 0.15}', '{"name": "Weapon Combo", "cost": 2, "damage_multiplier": 1.10, "debuff": {"stat": "enemy_technique", "value": -0.15}}', FALSE),
('Vin Jin', 'Allied', 440, 250, 2, 38, 42, 45, 1.60, 1.20, 1.25, '{"name": "Brutal Grip", "effect": "enemy_durability_debuff", "value": -0.10}', '{"name": "Grapple Break", "cost": 2, "damage_multiplier": 1.30, "debuff": {"stat": "enemy_durability", "value": -0.15}}', FALSE),
('Johan Seong', 'God Dog', 420, 260, 4, 60, 58, 34, 1.45, 1.55, 1.60, '{"name": "Copy Talent", "effect": "copy_highest_stat", "value": 0.10}', '{"name": "Copy Move", "cost": 3, "copy_str": 0.20, "damage_multiplier": 1.20}', FALSE),
('Jake Kim', 'Big Deal', 500, 245, 3, 36, 45, 60, 1.55, 1.15, 1.30, '{"name": "Leader Spirit", "effect": "squad_hp_buff", "value": 0.10}', '{"name": "Big Deal Command", "cost": 2, "squad_buff": {"stat": "str", "value": 0.15, "duration": 1}}', FALSE),
('Samuel Seo', 'Big Deal', 480, 270, 3, 40, 50, 55, 1.65, 1.20, 1.25, '{"name": "Crazy Mode", "effect": "str_buff_low_hp", "hp_threshold": 0.40, "value": 0.25}', '{"name": "Crazy Mode Hit", "cost": 2, "damage_multiplier": 1.40, "conditional_damage": {"hp_threshold": 0.50, "damage_multiplier": 1.70}}', FALSE),
('Eli Jang', 'Hostel', 460, 255, 3, 52, 50, 48, 1.50, 1.40, 1.35, '{"name": "Wild Instinct", "effect": "multi_buff", "speed": 0.15, "agility": 0.10}', '{"name": "Wild Hunt", "cost": 2, "damage_multiplier": 1.20, "buff": {"stat": "speed", "value": 0.20, "duration": 1}}', FALSE),
('Warren Chae', 'Hostel', 430, 235, 2, 40, 42, 52, 1.45, 1.20, 1.20, '{"name": "Protective Guard", "effect": "absorb_ally_damage", "value": 0.15}', '{"name": "Guard Counter", "cost": 1, "buff": {"stat": "incoming_damage", "value": -0.25, "duration": 1}, "counter_multiplier": 0.70}', FALSE),
('Jerry Kwon', 'Big Deal', 560, 275, 1, 28, 35, 65, 1.70, 1.00, 1.05, '{"name": "Iron Wall", "effect": "damage_reduction", "value": -0.20}', '{"name": "Heavy Smash", "cost": 1, "damage_multiplier": 1.35, "debuff": {"stat": "enemy_speed", "value": -0.10, "duration": 1}}', FALSE),
('Jason Yoon', 'Big Deal', 320, 180, 2, 55, 48, 32, 1.15, 1.50, 1.25, '{"name": "Fast Kick", "effect": "first_attack_chance", "chance": 0.25}', '{"name": "Fast Kick", "cost": 1, "effect": "first_attack", "damage_multiplier": 1.00}', FALSE),
('Jace Park', 'Burn Knuckles', 340, 160, 4, 38, 42, 38, 1.05, 1.20, 1.45, '{"name": "Tactical Mind", "effect": "squad_tech_buff", "value": 0.10}', '{"name": "Tactical Plan", "cost": 2, "squad_buff": {"stat": "technique", "value": 0.15, "duration": 1}}', FALSE),
('Duke Pyeon', 'J High', 300, 140, 5, 30, 35, 34, 1.00, 1.10, 1.35, '{"name": "Motivation", "effect": "heal_ally_once", "value": 40}', '{"name": "Motivation Song", "cost": 2, "effect": "heal_lowest_hp_ally", "value": 80}', FALSE),
('Jiho Park', 'None', 280, 150, 3, 36, 40, 30, 1.10, 1.15, 1.20, '{"name": "Dirty Trick", "effect": "enemy_str_debuff_chance", "chance": 0.20, "value": -0.15}', '{"name": "Dirty Trick", "cost": 2, "debuff": {"stat": "enemy_str", "value": -0.20, "duration": 1}}', FALSE),
('Gun Park', 'Campaign Boss', 600, 300, 5, 62, 68, 70, 1.80, 1.55, 1.70, '{"name": "Boss Mechanics"}', '{"name": "Boss Strike"}', TRUE);

-- Insert basic shop items
INSERT INTO shop_items (name, cost, description) VALUES
('Standard Character Pack', 1000, 'Roll for a random character'),
('Premium Character Pack', 2500, 'Roll for a random character with better odds');
