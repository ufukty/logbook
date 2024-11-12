package basics

type Boolean string

const (
	True  Boolean = "true"
	False Boolean = "false"
)

func (b *Boolean) FromRoute(src string) error {
	*b = Boolean(src)
	return b.Validate()
}

func (b Boolean) ToRoute() (string, error) {
	return string(b), nil
}
