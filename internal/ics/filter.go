package ics

import (
	"bytes"
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

func Filter(conf *config.Config, out *bytes.Buffer, in io.Reader) error { //nolint:gocognit,gocyclo,cyclop,funlen
	state := stateBegin
	stream := ics.NewCalendarStream(in)
	var wroteName bool
	serializeConfig := &ics.SerializationConfiguration{
		MaxLength:         75,
		PropertyMaxLength: 75,
		NewLine:           string(ics.NewLine),
	}

	for lineNum := 0; ; lineNum++ {
		l, err := stream.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if l == nil || len(*l) == 0 {
			continue
		}

		line, err := ics.ParseProperty(*l)
		if err != nil {
			return fmt.Errorf("%w %d: %w", ErrParsingLine, lineNum, err)
		}
		if line == nil {
			return fmt.Errorf("%w %d", ErrParsingCalendarLine, lineNum)
		}

		switch state {
		case stateBegin:
			switch line.IANAToken {
			case begin:
				if line.Value != vcalendar {
					return fmt.Errorf("%w: expected a vcalendar", ErrMalformed)
				}
				_ = line.SerializeTo(out, serializeConfig)
				state = stateProperties
			default:
				return fmt.Errorf("%w: expected begin", ErrMalformed)
			}
		case stateProperties:
			switch line.IANAToken {
			case begin:
				if conf.NewCalendarName != "" && !wroteName {
					nameProps := []*ics.BaseProperty{
						{IANAToken: string(ics.PropertyName), Value: conf.NewCalendarName},
						{IANAToken: string(ics.PropertyXWRCalName), Value: conf.NewCalendarName},
					}
					for _, prop := range nameProps {
						_ = prop.SerializeTo(out, serializeConfig)
					}
				}

				state = stateComponents
			case end:
				if line.Value != vcalendar {
					return fmt.Errorf("%w: expected end", ErrMalformed)
				}
				_ = line.SerializeTo(out, serializeConfig)
				state = stateEnd
			default:
				if slices.Contains(conf.CalendarFields, line.IANAToken) {
					switch ics.Property(line.IANAToken) {
					case ics.PropertyName, ics.PropertyXWRCalName:
						if conf.NewCalendarName != "" {
							wroteName = true
							line.Value = conf.NewCalendarName
						}
					}

					_ = line.SerializeTo(out, serializeConfig)
				}
			}
			if state != stateComponents {
				break
			}
			fallthrough
		case stateComponents:
			switch line.IANAToken {
			case end:
				if line.Value != vcalendar {
					return fmt.Errorf("%w: expected end", ErrMalformed)
				}
				_ = line.SerializeTo(out, serializeConfig)
				state = stateEnd
			case begin:
				co, err := ics.GeneralParseComponent(stream, line)
				if err != nil {
					return err
				}

				if co == nil || !slices.Contains(conf.Components, line.Value) {
					continue
				}

				if base := baseOf(co); base != nil {
					FilterComponent(conf, base)
				}

				_ = co.SerializeTo(out, serializeConfig)
			default:
				return fmt.Errorf("%w: expected begin or end", ErrMalformed)
			}
		case stateEnd:
			return fmt.Errorf("%w: unexpected end", ErrMalformed)
		default:
			return fmt.Errorf("%w: bad state", ErrMalformed)
		}
	}

	return nil
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
