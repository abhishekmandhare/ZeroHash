package models

// Trade struct is the standard representation of any trade in this application
type Trade struct {
	Currency string
	Price    float64
	Quantity float64
}

func (t Trade) GetKey() string {
	return t.Currency
}
