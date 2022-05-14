package vwap

import (
	"testing"

	"github.com/abhishekmandhare/zeroHash/internal/app/models"
	"github.com/abhishekmandhare/zeroHash/internal/app/stream"
	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {

	vwap := NewVwap(5, make(<-chan models.Trade), make(chan<- stream.StreamData))

	require.Equal(t, 25.1, vwap.calculate(models.Trade{Currency: "XXX", Price: 25.1, Quantity: 100}))                // 2510
	require.Equal(t, 25.133333333333333, vwap.calculate(models.Trade{Currency: "XXX", Price: 25.2, Quantity: 50}))   //1260
	require.Equal(t, 25.28, vwap.calculate(models.Trade{Currency: "XXX", Price: 25.5, Quantity: 100}))               //2550
	require.Equal(t, 25.589425867507888, vwap.calculate(models.Trade{Currency: "XXX", Price: 26.744, Quantity: 67})) //1791.848
	require.Equal(t, 25.446541436464088, vwap.calculate(models.Trade{Currency: "XXX", Price: 24.44, Quantity: 45}))  //1099.8
	require.Equal(t, 27.540363636363633, vwap.calculate(models.Trade{Currency: "XXX", Price: 30.11, Quantity: 200})) // 6022 12723.648 / 462

}
