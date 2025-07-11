-- Add response column to event_users table to track RSVP status
ALTER TABLE event_users ADD COLUMN response TEXT(20) DEFAULT 'going';

-- Update existing records to have 'going' as default response
UPDATE event_users SET response = 'going' WHERE response IS NULL;
