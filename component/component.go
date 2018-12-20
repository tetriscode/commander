package component

type Component interface {
	NewComponent(map[string]interface{})
	Start() error
	Stop(error)
}
