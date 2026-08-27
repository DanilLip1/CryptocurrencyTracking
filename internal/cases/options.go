package cases

type Mode int

const (
	ModeLatestPrices Mode = iota
	ModeMinPrices
	ModeMaxPrices
	ModePriceChangePercent
)

type Options struct {
	Mode Mode
}

type Option func(*Options)

func WithLatest() Option {
	return func(o *Options) {
		o.Mode = ModeLatestPrices
	}
}
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

func NewOptions(opts ...Option) *Options {
	options := &Options{
		Mode: ModeLatestPrices,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
