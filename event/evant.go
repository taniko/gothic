package event

type Type int

const (
	InvalidEvent Type = iota
	VersionEvent
	NameSpaceEvent
	StartSequenceEvent
	EndSequenceEvent
	ScalarEvent
	StartArrayEvent
	EndArrayEvent
	PropertyNameEvent
	PropertyTypeEvent
	DefineEvent
	NameEvent
	ServiceEvent
	ResourceEvent
	InputEvent
	OutputEvent
	ErrorsEvent
	CommentEvent
	OperationDefinitionEvent
	StructureDefinitionEvent
	ServiceVersionEvent
	ServiceOperationsEvent
	ServiceOperationEvent
	ServiceResourcesEvent
	OperationEvent
	OperationNameEvent
)

func New(t Type, v string) Event {
	return Event{
		t: t,
		v: v,
	}
}

type Event struct {
	t Type
	v string
}

func (s Event) Value() string {
	return s.v
}

func (s Event) Type() Type {
	return s.t
}
