-- Drop in reverse dependency order (FK constraints)
DROP TABLE IF EXISTS activity_logs;
DROP TABLE IF EXISTS share_links;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
