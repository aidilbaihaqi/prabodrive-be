-- Add document_name column to activity_logs for displaying file/folder names in activity feed
ALTER TABLE activity_logs ADD COLUMN IF NOT EXISTS document_name VARCHAR(255);
