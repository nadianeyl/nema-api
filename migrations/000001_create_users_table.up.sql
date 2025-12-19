CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY,
  name text NOT NULL DEFAULT '',
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  activated bool NOT NULL,
  email_notifications_enabled bool NOT NULL DEFAULT true,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1
);
