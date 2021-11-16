package mediator

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
	"github.com/ashwinath/financials/api/utils"
)

// ExpenseMediator is the mediator that calculates and stores the expenses.
type ExpenseMediator struct {
	expenseService  *service.ExpenseService
	expensesCSVFile string
}

// NewExpensesMediator creates a new ExpenseMediator
func NewExpensesMediator(
	expenseService *service.ExpenseService,
	expensesCSVFile string,
) *ExpenseMediator {
	return &ExpenseMediator{
		expenseService:  expenseService,
		expensesCSVFile: expensesCSVFile,
	}
}

func (m *ExpenseMediator) parseCSV() ([]*models.Expense, error) {
	records, err := utils.ReadCSV(m.expensesCSVFile)
	if err != nil {
		return nil, err
	}

	var expenses []*models.Expense
	headers := records[0]
	for recordNum := 1; recordNum < len(records); recordNum++ {
		expense := &models.Expense{}
		for i, value := range records[recordNum] {
			switch headers[i] {
			case "date":
				layout := "2006-01-02T15:04:05.000Z"
				str := fmt.Sprintf("%sT08:00:00.000Z", value)
				t, err := time.Parse(layout, str)
				if err != nil {
					return nil, err
				}
				expense.TransactionDate = t
			case "type":
				expense.Type = value
			case "amount":
				if v, err := strconv.ParseFloat(value, 64); err == nil {
					expense.Amount = v
				} else {
					return nil, err
				}
			}
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// ProcessExpenses reads the csvs and stores them
func (m *ExpenseMediator) ProcessExpenses() {
	expenses, err := m.parseCSV()
	if err != nil {
		log.Printf("Could not parse CSV: %s", err)
		return
	}

	err = m.expenseService.TruncateTable()
	if err != nil {
		log.Printf("Could not truncate expenses table: %s", err)
		return
	}

	err = m.expenseService.BulkAdd(expenses)
	if err != nil {
		log.Printf("Could not bulk add expenses: %s", err)
		return
	}
}
