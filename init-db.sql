-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the database if it doesn't exist
-- (This is handled by POSTGRES_DB environment variable) 