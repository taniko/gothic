package event

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type test struct {
		name  string
		input string
		want  []Event
	}
	tests := []test{
		{
			name:  "version",
			input: `$version: "2"`,
			want:  []Event{New(VersionEvent, "2")},
		}, {
			name:  "namespace",
			input: "namespace com.example",
			want:  []Event{New(NameSpaceEvent, "com.example")},
		}, {
			name:  "service",
			input: "service Weather",
			want:  []Event{New(ServiceEvent, "service"), New(NameEvent, "Weather")},
		}, {
			name:  "service with start sequence",
			input: "service Weather{",
			want:  []Event{New(ServiceEvent, "service"), New(NameEvent, "Weather"), New(StartSequenceEvent, "{")},
		}, {
			name:  "service with start sequence",
			input: "service Weather {",
			want:  []Event{New(ServiceEvent, "service"), New(NameEvent, "Weather"), New(StartSequenceEvent, "{")},
		}, {
			name:  "service version",
			input: `version: "2006-03-01"`,
			want: []Event{
				New(ServiceVersionEvent, "2006-03-01"),
			},
		}, {
			name:  "operation",
			input: "operation GetCity{",
			want: []Event{
				New(OperationDefinitionEvent, "operation"),
				New(NameEvent, "GetCity"),
				New(StartSequenceEvent, "{"),
			},
		}, {
			name:  "operation",
			input: "operation GetCity {",
			want: []Event{
				New(OperationDefinitionEvent, "operation"),
				New(NameEvent, "GetCity"),
				New(StartSequenceEvent, "{"),
			},
		}, {
			name:  "end sequence",
			input: "}",
			want:  []Event{New(EndSequenceEvent, "}")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := Parse(bufio.NewScanner(strings.NewReader(tt.input)))
			assert.NoError(t, err)
			for i, e := range tt.want {
				t.Log(i, len(tt.want))
				v := events.Event()
				assert.Equal(t, e.Type(), v.Type())
				assert.Equal(t, e.Value(), v.Value())
				assert.Equal(t, i+1 < len(tt.want), events.Next())
			}
		})
	}
}
