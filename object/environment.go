package object

type Environment struct {
	data   map[string]Object
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		data: make(map[string]Object),
	}
}

func (e *Environment) Get(key string) (Object, bool) {
	v, ok := e.data[key]
	if !ok && e.parent != nil {
		return e.parent.Get(key)
	}
	return v, ok
}
func (e *Environment) Set(key string, value Object) {
	e.data[key] = value
}
func (e *Environment) NewScope() *Environment {
	env := NewEnvironment()
	env.parent = e
	return env
}
