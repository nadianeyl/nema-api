CREATE TYPE account_class AS ENUM ('cce', 'investment', 'liability');

CREATE TABLE IF NOT EXISTS accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,
  name text NOT NULL,
  class account_class NOT NULL,
  currency_code char(3) NOT NULL,
  balance decimal(15, 2) NOT NULL DEFAULT 0,
  is_budgeted bool NOT NULL DEFAULT FALSE,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts (user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_user_id_class ON accounts (user_id, class);
