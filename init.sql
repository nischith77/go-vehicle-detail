-- Create models table for storing vehicle data
CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY,
    make_id INTEGER NOT NULL,
    make VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on make_id for better query performance
CREATE INDEX IF NOT EXISTS idx_models_make_id ON models(make_id);

-- Create index on make for better query performance
CREATE INDEX IF NOT EXISTS idx_models_make ON models(make);


-- Display confirmation
SELECT 'Models table created successfully!' as status;