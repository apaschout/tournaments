package tournaments

type Phase string

type Format interface {
	Type() string
	Init() Phase
}

type CubePauper struct {
}

func (cp *CubePauper) Type() string {
	return "cube-pauper"
}

func (cp *CubePauper) Init() Phase {
	return "registration"
}
