package coffeeco

import (
	"github.com/govalues/decimal"
)

// Product is value object
// - no identity? if you change the name than it is different product
// - can be treated as immutable
// - can be compared only by values
type Product struct {
	ItemName  string
	BasePrice decimal.Decimal
}
