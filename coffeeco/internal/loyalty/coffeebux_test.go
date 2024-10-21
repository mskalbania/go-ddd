package loyalty

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddStampResultsInNewFreeDrink(t *testing.T) {
	bux := CoffeeBux{
		FreeDrinksAvailable:                   1,
		RemainingDrinkPurchasesUntilFreeDrink: 1,
	}

	newBux := bux.AddStamp()

	require.Equal(t, bux.FreeDrinksAvailable+1, newBux.FreeDrinksAvailable)
	require.Equal(t, 10, newBux.RemainingDrinkPurchasesUntilFreeDrink)
}

func TestAddStampSuccessful(t *testing.T) {
	bux := CoffeeBux{
		RemainingDrinkPurchasesUntilFreeDrink: 5,
	}

	newBux := bux.AddStamp()

	require.Equal(t, bux.FreeDrinksAvailable, newBux.FreeDrinksAvailable)
	require.Equal(t, bux.RemainingDrinkPurchasesUntilFreeDrink-1, newBux.RemainingDrinkPurchasesUntilFreeDrink)
}
