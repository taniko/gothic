package node

type Resource struct {
	Identifiers string
	Create      *Operation
	Put         *Operation
	Read        *Operation
	Delete      *Operation
	List        *Operation
}
