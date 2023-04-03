package types

// Resource is the actual representation of a resource.
type Resource struct {
	Identifier Identifier
	State      string
	TypeName   string
}
