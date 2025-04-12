package object

// Environment represents the running interpreter Environment
// including bound values to identifiers.
type Environment struct {
	store map[string]Object // store is the persistent store for identifiers and their values.
}

// NewEnvironment creates a new environment.
//
// Returns:
//   - *Environment: the newly created environment.
func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store}
}

// Get looks for an identifier in the environment store.
//
// Parameters:
//   - name: The name of the identifier to get.
//
// Returns:
//   - Object: The object if it's found.
//   - bool: True when the object was found, otherwise false.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set stores the name and value of the identifier in the environment store.
//
// Parameters:
//   - name: The name of the user defined identifier.
//   - val: The object to store.
//
// Returns:
//   - Object: The input object after saving.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
