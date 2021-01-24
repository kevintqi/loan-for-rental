package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"math"
)

const (
	StandardExpenseFactor = 0.35
	CashFlowWeight        = 1.2
	ThirtyYears           = 360
)

type Input struct {
	AskPrice     float64 `json:"ask_price"`
	Income       float64 `json:"income"`
	Expense      float64 `json:"expense"`
	InterestRate float64 `json:"interest_rate"`
}

type Result struct {
	InputData         Input   `json:"input:"`
	EffectiveCashFlow float64 `json:"effective_cash_flow"`
	LoanAmount        float64 `json:"loan_amount:"`
	Delta             float64 `json:"delta"`
}

func effectiveCashFlow(event Input) float64 {
	expense := event.Income * StandardExpenseFactor
	if event.Expense > expense {
		expense = event.Expense
	}
	cashFlow := (event.Income - expense) / CashFlowWeight
	return cashFlow
}

func monthlyRate(event Input) float64 {
	return (event.InterestRate / 100.0) / 12.0
}

func loanCalc(event Input, cashFlow float64) float64 {
	r := monthlyRate(event)
	ftr := (1.0 - math.Pow(1+r, -ThirtyYears)) / r
	loanAmount := cashFlow * ftr
	return loanAmount
}

func HandleLambdaEvent(ctx context.Context, event Input) (Result, error) {
	cashFlow := effectiveCashFlow(event)
	loanAmount := loanCalc(event, cashFlow)
	delta := event.AskPrice - loanAmount
	return Result{InputData: event, EffectiveCashFlow: cashFlow, LoanAmount: loanAmount, Delta: delta}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
