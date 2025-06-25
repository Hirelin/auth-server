-- Active: 1750847222187@@db.cgsbvqczhzwqtkrucmgd.supabase.co@5432@hirelin-db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE DATABASE "hirelin-db";

-- createTable
-- User table
CREATE TABLE "User" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    image TEXT DEFAULT '/images/profile-default.jpg',
    username TEXT UNIQUE NOT NULL,

    email_verified TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Account table
CREATE TABLE "Account" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID REFERENCES "User"(id) ON DELETE CASCADE,
    type TEXT CHECK (type IN ('oauth', 'email', 'credentials')) DEFAULT 'oauth',
    provider TEXT NOT NULL,
    provider_account_id TEXT NOT NULL,

    access_token TEXT,
    refresh_token TEXT,

    expires_at TIMESTAMP,
    
    token_type TEXT,
    scope TEXT,
    id_token TEXT,
    session_state TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- VerificationToken table
CREATE TABLE "VerificationToken" (
    identifier TEXT NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,

    PRIMARY KEY (identifier, token)
);

-- Sessions table
CREATE TABLE "Session" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    session_token TEXT UNIQUE NOT NULL,
    refresh_token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    user_id UUID REFERENCES "User"(id) ON DELETE CASCADE
);

-- Function to set updatedAt timestamps
CREATE OR REPLACE FUNCTION set_current_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers to set updated_at timestamps for all tables
CREATE TRIGGER updated_at_user_trigger
BEFORE UPDATE ON "User"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();

CREATE TRIGGER updated_at_session_trigger
BEFORE UPDATE ON "Session"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();