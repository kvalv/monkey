package object

type Environment struct {
	data map[string]Object
}

func New() *Environment {
	return &Environment{
		data: make(map[string]Object),
	}
}

func (e *Environment) Get(key string) (Object, bool) {
	v, ok := e.data[key]
	return v, ok
}
func (e *Environment) Set(key string, value Object) {
	e.data[key] = value
}
