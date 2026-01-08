CREATE TABLE budget_items (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  budget_id uuid NOT NULL REFERENCES budgets ON DELETE CASCADE,
  category_id uuid NOT NULL REFERENCES categories ON DELETE RESTRICT,
  limit_amount decimal(15, 2) NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1,

  CONSTRAINT cc_budget_items_limit_amount CHECK (limit_amount > 0),
  CONSTRAINT uc_budget_items_budget_category UNIQUE (budget_id, category_id)
);

CREATE INDEX idx_budget_items_budget_id ON budget_items (budget_id);
CREATE INDEX idx_budget_items_category_id ON budget_items (category_id);
