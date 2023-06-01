package alilog

type optionType string

type (
	optionValue struct {
		Value interface{}
		Type  optionType
	}

	// Option HTTP option
	Option func(map[string]optionValue) error
)
