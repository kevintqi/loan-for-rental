package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"math"
	"strconv"
)

const (
	StandardExpenseFactor = 0.35
	CashFlowWeight        = 1.2
	ThirtyYears           = 360
	AskPriceKey           = "ask_price"
	IncomeKey             = "income"
	ExpenseKey            = "expense"
	InterestRateKey       = "interest_rate"
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

func buildInput(query map[string]string) (Input, error) {
	askPrice, err := strconv.ParseFloat(query[AskPriceKey], 64)
	if err != nil {
		return Input{}, err
	}
	income, err := strconv.ParseFloat(query[IncomeKey], 64)
	if err != nil {
		return Input{}, err
	}
	expense, err := strconv.ParseFloat(query[ExpenseKey], 64)
	if err != nil {
		return Input{}, err
	}
	interestRate, err := strconv.ParseFloat(query[InterestRateKey], 64)
	if err != nil {
		return Input{}, err
	}
	return Input{AskPrice: askPrice, Income: income, Expense: expense, InterestRate: interestRate}, nil
}

func generateResult(input Input) Result {
	cashFlow := effectiveCashFlow(input)
	loanAmount := loanCalc(input, cashFlow)
	delta := input.AskPrice - loanAmount
	return Result{InputData: input, EffectiveCashFlow: cashFlow, LoanAmount: loanAmount, Delta: delta}
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input, err := buildInput(request.QueryStringParameters)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}
	result := generateResult(input)
	body, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
