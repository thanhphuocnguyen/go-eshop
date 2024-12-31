package payment

type PaymentParam struct {
	Amount float64
	Email  string
	Source string
}
type PaymentResult struct {
	TransactionID string
	Amount        float64
	Currency      string
	Source        string
	Email         string
}

type PaymentService interface {
	MakePayment(param PaymentParam) (PaymentResult, error)
	CancelPayment(transactionID string) error
}
