package payment

type Means string

const (
	CARD      = Means("card")
	CASH      = Means("cash")
	COFFEEBUX = Means("coffeebux")
)

// CardDetails is simplification over cards, assumptions is we will get the token and charge it
type CardDetails struct {
	cardToken string
}