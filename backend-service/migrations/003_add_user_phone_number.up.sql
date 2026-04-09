-- Add phone_number to sys_user table
ALTER TABLE sys_user ADD COLUMN IF NOT EXISTS phone_number VARCHAR(50);

-- Create index for phone_number
CREATE INDEX IF NOT EXISTS idx_sys_user_phone_number ON sys_user(phone_number);
