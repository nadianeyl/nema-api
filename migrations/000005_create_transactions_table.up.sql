CREATE TABLE IF NOT EXISTS transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,
  type transaction_type NOT NULL,
  category_id uuid NOT NULL REFERENCES categories ON DELETE RESTRICT,
  amount decimal(15, 2) NOT NULL,
  date timestamp with time zone NOT NULL DEFAULT NOW(),
  title text,
  notes text,
  from_account_id uuid REFERENCES accounts ON DELETE RESTRICT,
  to_account_id uuid REFERENCES accounts ON DELETE RESTRICT,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1,

  CONSTRAINT cc_transactions_amount CHECK (amount > 0),
  CONSTRAINT cc_transactions_expense_requires_from_account CHECK (type != 'expense' OR from_account_id IS NOT NULL),
  CONSTRAINT cc_transactions_income_requires_to_account CHECK (type != 'income' OR to_account_id IS NOT NULL),
  CONSTRAINT cc_transactions_transfer_requires_both_accounts CHECK (type != 'transfer' OR (from_account_id IS NOT NULL AND to_account_id IS NOT NULL)),
  CONSTRAINT cc_transactions_transfer_different_accounts CHECK (type != 'transfer' OR from_account_id != to_account_id)
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id_date ON transactions (user_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_type ON transactions (user_id, type);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_category_id ON transactions (user_id, category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_type_date ON transactions (user_id, type, date DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_from_account_id ON transactions (from_account_id) WHERE from_account_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_transactions_to_account_id ON transactions (to_account_id) WHERE to_account_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_transactions_title ON transactions USING GIN (to_tsvector('simple', title));
