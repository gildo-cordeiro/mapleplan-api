package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_NewUser(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)

	email := "gildo.duarte"
	password := "123456"
	firstName := "Gildo"
	lastName := "Duarte"

	transaction := Transaction{
		Model:           gorm.Model{ID: 1},
		UserID:          1,
		Amount:          decimal.NewFromFloat(123.45),
		Currency:        "USD",
		Type:            "deposit",
		Description:     "Initial deposit",
		TransactionDate: time.Date(2018, time.April, 26, 1, 2, 3, 0, time.UTC),
	}

	task := Task{
		Model:       gorm.Model{ID: 2},
		UserID:      1,
		Title:       "Finish report",
		Description: "monthly finance report",
		DueDate:     time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
		IsCompleted: false,
	}

	goal := Goal{
		Model:         gorm.Model{ID: 3},
		UserID:        1,
		Name:          "Vacation",
		TargetAmount:  decimal.NewFromInt(2000),
		CurrentAmount: decimal.NewFromFloat(150.75),
		DueDate:       time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC),
	}

	user, _ := NewUser(email, password, firstName, lastName, &transaction, &task, &goal)

	assertions.NotNilf(user, "NewUser must not be nil")
	assertions.Equal(email, user.Email)
	assertions.Equal(password, user.PasswordHash)
	assertions.Equal(firstName, user.FirstName)
	assertions.Equal(lastName, user.LastName)

	assertions.Equal(transaction.ID, user.Transaction.ID)
	assertions.Equal(transaction.Amount.String(), user.Transaction.Amount.String())
	assertions.Equal(transaction.Currency, user.Transaction.Currency)
	assertions.Equal(transaction.Type, user.Transaction.Type)

	assertions.Equal(task.ID, user.Task.ID)
	assertions.Equal(task.Title, user.Task.Title)
	assertions.Equal(task.IsCompleted, user.Task.IsCompleted)

	assertions.Equal(goal.ID, user.Goal.ID)
	assertions.Equal(goal.Name, user.Goal.Name)
	assertions.Equal(goal.TargetAmount.String(), user.Goal.TargetAmount.String())

}
