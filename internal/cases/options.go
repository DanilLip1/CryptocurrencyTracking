package cases

type Mode int

const (
	_ Mode = iota
	ModeMinPrices
	ModeMaxPrices
	ModePriceChangePercent
)

type Options struct {
	Mode Mode
}

type Option func(*Options)

func WithMin() Option {
	return func(o *Options) {
		o.Mode = ModeMinPrices
	}
}
func WithMax() Option {
	return func(o *Options) {
		o.Mode = ModeMaxPrices
	}
}
func WithPercent() Option {
	return func(o *Options) {
		o.Mode = ModePriceChangePercent
	}
}
