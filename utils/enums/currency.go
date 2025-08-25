package enums

type Currency int64

const (
	UNDEFINED Currency = iota
	EUR
	USD
)

func (c Currency) String() string {
	switch c {
	case EUR:
		return "EUR"
	case USD:
		return "USD"
	}
	return "UNDEFINED"
}
