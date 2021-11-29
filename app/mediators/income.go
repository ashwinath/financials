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

// IncomeMediator is the mediator that calculates and stores the incomes.
type IncomeMediator struct {
	incomeService  *service.IncomeService
	incomesCSVFile string
}

// NewIncomeMediator creates a new IncomeMediator
func NewIncomeMediator(
	incomeService *service.IncomeService,
	incomesCSVFile string,
) *IncomeMediator {
	return &IncomeMediator{
		incomeService:  incomeService,
		incomesCSVFile: incomesCSVFile,
	}
}

func (m *IncomeMediator) parseCSV() ([]*models.Income, error) {
	records, err := utils.ReadCSV(m.incomesCSVFile)
	if err != nil {
		return nil, err
	}

	var incomes []*models.Income
	headers := records[0]
	for recordNum := 1; recordNum < len(records); recordNum++ {
		income := &models.Income{}
		for i, value := range records[recordNum] {
			switch headers[i] {
			case "date":
				layout := "2006-01-02T15:04:05.000Z"
				str := fmt.Sprintf("%sT08:00:00.000Z", value)
				t, err := time.Parse(layout, str)
				if err != nil {
					return nil, err
				}
				income.TransactionDate = t
			case "type":
				income.Type = value
			case "amount":
				if v, err := strconv.ParseFloat(value, 64); err == nil {
					income.Amount = v
				} else {
					return nil, err
				}
			}
		}
		incomes = append(incomes, income)
	}

	return incomes, nil
}

// ProcessIncome reads the csvs and stores them
func (m *IncomeMediator) ProcessIncome() error {
	log.Printf("Updating incomes.")
	incomes, err := m.parseCSV()
	if err != nil {
		return err
	}

	err = m.incomeService.TruncateTable()
	if err != nil {
		return err
	}

	err = m.incomeService.BulkAdd(incomes)
	if err != nil {
		return err
	}

	return nil
}
