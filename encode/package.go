package encode

import (
	"errors"

	"github.com/taniko/gothic/event"
	"github.com/taniko/gothic/node"
)

type Package struct {
	NameSpace string
	Services  []*node.Service
}

func NewPackage(scanner *event.Scanner) (*Package, error) {
	p := Package{}
	for scanner.Next() {
		e := scanner.Event()
		switch e.Type() {
		case event.NameSpaceEvent:
			p.NameSpace = e.Value()
		case event.ServiceEvent:
			service, err := scanService(scanner)
			if err != nil {
				return nil, err
			}
			p.Services = append(p.Services, service)
		}
	}
	return &p, nil
}

func scanService(scanner *event.Scanner) (*node.Service, error) {
	service := node.Service{}
	for scanner.Next() {
		e := scanner.Event()
		switch e.Type() {
		case event.NameEvent:
			service.Name = e.Value()
		case event.ServiceVersionEvent:
			service.Version = e.Value()
		case event.StartSequenceEvent:
			continue
		case event.EndSequenceEvent:
			return &service, nil
		case event.VersionEvent:
			service.Version = e.Value()
		case event.ServiceOperationsEvent:
			operations, err := scanServiceOperation(scanner)
			if err != nil {
				return nil, err
			}
			service.Operations = operations
		}
	}
	return nil, errors.New("not find sequence end character")
}

func scanServiceOperation(scanner *event.Scanner) ([]string, error) {
	var operations []string
	for scanner.Next() {
		e := scanner.Event()
		switch e.Type() {
		case event.StartArrayEvent:
			continue
		case event.EndArrayEvent:
			return operations, nil
		case event.ServiceOperationEvent:
			operations = append(operations, e.Value())
		}
	}
	return nil, errors.New("not find array end character")
}
