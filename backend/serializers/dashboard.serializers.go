package serializers

type DashboardResponse struct {
	Summary            DashboardSummary      `json:"summary"`
	Graphs             DashboardGraphs       `json:"graphs"`
	RecentTransactions []TransactionResponse `json:"recentTransactions"`
}

type DashboardSummary struct {
	TotalBalance    float64 `json:"totalBalance"`
	MonthlyIncome   float64 `json:"monthlyIncome"`
	MonthlyExpenses float64 `json:"monthlyExpenses"`
	MonthlySavings  float64 `json:"monthlySavings"`
}

type DashboardGraphs struct {
	IncomeExpenseAggregate IncomeExpenseAggregate `json:"incomeExpenseAggregate"`
	SpendingByCategory     []CategoryAggregate    `json:"spendingByCategory"`
	EarningByCategory      []CategoryAggregate    `json:"earningByCategory"`
}

type CategoryAggregate struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

type IncomeExpenseAggregate struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}
