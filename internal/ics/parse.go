package ics

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"gabe565.com/ics-redact-proxy/internal/config"
	ics "github.com/arran4/golang-ical"
)

type parseState uint8

const (
	stateBegin parseState = iota
	stateProperties
	stateComponents
	stateEnd
)

const (
	begin     string = "BEGIN"
	end       string = "END"
	vcalendar string = "VCALENDAR"
)

var (
	ErrParsingLine         = errors.New("parsing line")
	ErrParsingCalendarLine = errors.New("parsing calendar line")
	ErrMalformed           = errors.New("malformed calendar")
)

func ParseAndFilter(conf *config.Config, r io.Reader) (*ics.Calendar, error) { //nolint:gocognit,funlen
	state := stateBegin
	cal := &ics.Calendar{}
	stream := ics.NewCalendarStream(r)

	for lineNum := 0; ; lineNum++ {
		l, err := stream.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return cal, err
		}
		if l == nil || len(*l) == 0 {
			continue
		}

		line, err := ics.ParseProperty(*l)
		if err != nil {
			return nil, fmt.Errorf("%w %d: %w", ErrParsingLine, lineNum, err)
		}
		if line == nil {
			return nil, fmt.Errorf("%w %d", ErrParsingCalendarLine, lineNum)
		}

		switch state {
		case stateBegin:
			switch line.IANAToken {
			case begin:
				if line.Value != vcalendar {
					return nil, fmt.Errorf("%w: expected a vcalendar", ErrMalformed)
				}
				state = stateProperties
			default:
				return nil, fmt.Errorf("%w: expected begin", ErrMalformed)
			}
		case stateProperties:
			switch line.IANAToken {
			case begin:
				state = stateComponents
			case end:
				if line.Value != vcalendar {
					return nil, fmt.Errorf("%w: expected end", ErrMalformed)
				}
				state = stateEnd
			default:
				cal.CalendarProperties = append(cal.CalendarProperties, ics.CalendarProperty{BaseProperty: *line})
			}
			if state != stateComponents {
				break
			}
			fallthrough
		case stateComponents:
			switch line.IANAToken {
			case end:
				if line.Value != vcalendar {
					return nil, fmt.Errorf("%w: expected end", ErrMalformed)
				}
				state = stateEnd
			case begin:
				co, err := ics.GeneralParseComponent(stream, line)
				if err != nil {
					return nil, err
				}

				if co == nil || !slices.Contains(conf.Components, line.Value) {
					continue
				}

				if base := baseOf(co); base != nil {
					FilterComponent(conf, base)
				}
				cal.Components = append(cal.Components, co)
			default:
				return nil, fmt.Errorf("%w: expected begin or end", ErrMalformed)
			}
		case stateEnd:
			return nil, fmt.Errorf("%w: unexpected end", ErrMalformed)
		default:
			return nil, fmt.Errorf("%w: bad state", ErrMalformed)
		}
	}

	if conf.NewCalendarName != "" {
		cal.SetName(conf.NewCalendarName)
	}

	return cal, nil
}

func baseOf(component any) *ics.ComponentBase {
	switch t := component.(type) {
	case *ics.VEvent:
		return &t.ComponentBase
	case *ics.VTodo:
		return &t.ComponentBase
	case *ics.VJournal:
		return &t.ComponentBase
	case *ics.VBusy:
		return &t.ComponentBase
	default:
		return nil
	}
}
