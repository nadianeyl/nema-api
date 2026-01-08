CREATE TABLE budgets (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,
  name text NOT NULL,
  available_budget decimal(15, 2) NOT NULL,
  start_date date NOT NULL,
  end_date date NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1,
  
  CONSTRAINT cc_budgets_available_budget CHECK (available_budget >= 0),
  CONSTRAINT cc_budgets_date_range_valid CHECK (end_date > start_date),
  CONSTRAINT uc_budgets_user_date_range UNIQUE (user_id, start_date, end_date)
);

CREATE INDEX idx_budgets_user_id ON budgets (user_id);
CREATE INDEX idx_budgets_user_id_name ON budgets (user_id, name);
CREATE INDEX idx_budgets_user_id_date_range ON budgets (user_id, start_date, end_date);
