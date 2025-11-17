-- Initialize PostgreSQL database schema for CodexGigantus
-- This script creates the necessary tables for storing code files

-- Create code_files table
CREATE TABLE IF NOT EXISTS code_files (
    id SERIAL PRIMARY KEY,
    file_path VARCHAR(1024) NOT NULL,
    content TEXT NOT NULL,
    file_type VARCHAR(50),
    file_size INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on file_path for faster lookups
CREATE INDEX IF NOT EXISTS idx_file_path ON code_files(file_path);

-- Create index on file_type
CREATE INDEX IF NOT EXISTS idx_file_type ON code_files(file_type);

-- Insert sample data (optional)
INSERT INTO code_files (file_path, content, file_type, file_size) VALUES
    ('main.go', 'package main\n\nfunc main() {\n    println("Hello, World!")\n}', 'go', 67),
    ('utils.py', 'def hello():\n    print("Hello from Python")', 'py', 45),
    ('app.js', 'console.log("Hello from JavaScript");', 'js', 39)
ON CONFLICT DO NOTHING;

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_code_files_updated_at BEFORE UPDATE
    ON code_files FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Grant permissions (adjust as needed)
GRANT ALL PRIVILEGES ON TABLE code_files TO codex;
GRANT USAGE, SELECT ON SEQUENCE code_files_id_seq TO codex;
