package gogrep

type Errors []error

// implementation of error interface
func (e Errors) Error() string {
	var err string
	for i, s := range e {
		err += s.Error()
		if i > 0 {
			err += "\n"
		}
	}
	return  err
}

func (e Errors) String() string {
	return e.Error()
}

func (e Errors) HasErrors() bool {
	return e.Len() > 0
}

func (e Errors) Len() int {
	return len(e)
}

func (e Errors) Walk(callable func(err error)) {
	for _, err := range e {
		callable(err)
	}
}

func (e Errors) add(err error) {
	e = append(e, err)
}