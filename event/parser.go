package event

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	startSequence          = regexp.MustCompile(`\{$`)
	endSequence            = regexp.MustCompile(`^}$`)
	smithyVersionEvent     = regexp.MustCompile(`^\$version:\s"(.+)"$`)
	defineEventRegex       = regexp.MustCompile(`^(service|resource|operation|structure)\s([a-zA-Z0-9]+)(\s*{?)`)
	namespaceRegex         = regexp.MustCompile(`^namespace\s([a-zA-Z.-_]+)`)
	propertyRegex          = regexp.MustCompile(`^([a-zA-Z]+):\s?([a-zA-Z-0-9-"]+)`)
	serviceOperationsRegex = regexp.MustCompile(`operations:\s?\[(.+)]`)
)

type Scanner struct {
	index  int
	events []Event
}

func NewScanner(events []Event) Scanner {
	return Scanner{
		events: events,
	}
}

func (s *Scanner) Next() bool {
	ok := len(s.events) > s.index+1
	if ok {
		s.index += 1
	}
	return ok
}

func (s *Scanner) Event() Event {
	return s.events[s.index]
}

func (s *Scanner) Reset() {
	s.index = 0
}

func Parse(scanner *bufio.Scanner) (*Scanner, error) {
	var events []Event
	for scanner.Scan() {
		text := strings.TrimSpace(strings.TrimSuffix(scanner.Text(), "\n"))
		if m := smithyVersionEvent.FindStringSubmatch(text); len(m) == 2 {
			events = append(events, New(VersionEvent, m[1]))
		} else if m := defineEventRegex.FindStringSubmatch(text); len(m) >= 3 {
			t := InvalidEvent
			switch m[1] {
			case "service":
				t = ServiceEvent
			case "resource":
				t = ResourceEvent
			case "operation":
				t = OperationDefinitionEvent
			case "structure":
				t = StructureDefinitionEvent
			}
			if t != InvalidEvent {
				events = append(events,
					New(t, m[1]),
					New(NameEvent, m[2]),
				)
				if strings.TrimSpace(m[3]) == "{" {
					events = append(events, New(StartSequenceEvent, "{"))
				}
			}
		} else if startSequence.MatchString(text) {
			events = append(events, New(StartSequenceEvent, "{"))
		} else if endSequence.MatchString(text) {
			events = append(events, New(EndSequenceEvent, "}"))
		} else if m := propertyRegex.FindStringSubmatch(text); len(m) > 2 {
			switch m[1] {
			case "version":
				events = append(events, New(ServiceVersionEvent, strings.Trim(m[2], `"`)))
			case "input":
				events = append(events, New(InputEvent, m[2]))
			case "output":
				events = append(events, New(OutputEvent, m[2]))
			case "create", "put", "read", "update", "delete", "list":
				events = append(events, New(OperationEvent, m[1]), New(OperationNameEvent, m[2]))
			default:
				events = append(events, New(PropertyNameEvent, m[1]), New(PropertyTypeEvent, m[2]))
			}
		} else if m := namespaceRegex.FindStringSubmatch(text); len(m) > 1 {
			events = append(events, New(NameSpaceEvent, m[1]))
		} else if m := serviceOperationsRegex.FindStringSubmatch(text); len(m) > 1 {
			events = append(events,
				New(ServiceOperationsEvent, "operations"),
				New(StartArrayEvent, "start-operations"),
			)
			operations := strings.Split(m[1], ",")
			for _, operation := range operations {
				events = append(events, New(ServiceOperationEvent, strings.TrimSpace(operation)))
			}
			events = append(events, New(EndArrayEvent, "end-operations"))
		}
	}
	s := NewScanner(events)
	return &s, nil
}
