package domain

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

func (t TransactionType) String() string {
	return string(t)
}

func GetTransactionTypes() []TransactionType {
	return []TransactionType{TransactionTypeIncome, TransactionTypeExpense, TransactionTypeTransfer}
}
