CREATE TYPE transaction_type AS ENUM ('income', 'expense', 'transfer');

CREATE TABLE IF NOT EXISTS categories (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES users ON DELETE CASCADE,
  name text NOT NULL,
  transaction_type transaction_type NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1,  

  CONSTRAINT uc_categories_user_name_type UNIQUE NULLS NOT DISTINCT (user_id, name, transaction_type)
);

CREATE INDEX IF NOT EXISTS idx_categories_user_id_transaction_type ON categories (user_id, transaction_type);
CREATE INDEX IF NOT EXISTS idx_categories_transaction_type ON categories (transaction_type);
