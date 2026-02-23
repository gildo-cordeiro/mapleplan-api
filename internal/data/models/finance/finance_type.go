package finance

type FinanceType string

const (
	Income   FinanceType = "income"
	Expense  FinanceType = "expense"
	Transfer FinanceType = "transfer"
)

func IsValidType(financeType FinanceType) bool {
	return financeType == Income || financeType == Expense || financeType == Transfer
}

type Category string

const (
	Salary         Category = "salary"
	Investment     Category = "investment"
	Savings        Category = "savings"
	Housing        Category = "housing"
	Transportation Category = "transportation"
	Food           Category = "food"
	Healthcare     Category = "healthcare"
	Education      Category = "education"
	Utilities      Category = "utilities"
	Entertainment  Category = "entertainment"
	Other          Category = "other"
)

func IsValidCategory(category Category) bool {
	validCategories := map[Category]bool{
		Salary:         true,
		Investment:     true,
		Savings:        true,
		Housing:        true,
		Transportation: true,
		Food:           true,
		Healthcare:     true,
		Education:      true,
		Utilities:      true,
		Entertainment:  true,
		Other:          true,
	}
	return validCategories[category]
}
