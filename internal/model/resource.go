package model

type ResourceChange struct {
	Address string                 // resource address e.g. module.main.aws_instance.example
	Actions []string               // plan actions e.g. ["update"]
	Before  map[string]interface{} // attributes before change
	After   map[string]interface{} // attributes after change
}
