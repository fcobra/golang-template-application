CREATE TABLE IF NOT EXISTS catalog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    disabled BOOLEAN NOT NULL DEFAULT FALSE
);

-- Seed data
INSERT INTO catalog (title, description, disabled) VALUES
('Item 1', 'Description for item 1', FALSE),
('Item 2', 'Description for item 2', FALSE),
('Item 3', 'A disabled item', TRUE),
('Item 4', 'Another active item', FALSE);
