-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS predictions;
DROP TABLE IF EXISTS awards;
DROP TABLE IF EXISTS nominations;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;

-- Optionally drop the UUID extension (commented out to be safe)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
