package ics

import (
	"crypto/sha256"
	"encoding/hex"
	"slices"

	"gabe565.com/ics-redact-proxy/internal/config"
	ics "github.com/arran4/golang-ical"
)

func FilterEvent(conf *config.Config, event *ics.VEvent) {
	event.Properties = slices.DeleteFunc(event.Properties, func(property ics.IANAProperty) bool {
		return !slices.Contains(conf.EventAllowFields, property.IANAToken)
	})
	event.Properties = slices.Clip(event.Properties)

	if conf.NewEventSummary != "" {
		event.SetSummary(conf.NewEventSummary)
	}

	if conf.HashUID {
		if uid := event.GetProperty(ics.ComponentPropertyUniqueId); uid != nil {
			checksum := sha256.Sum256([]byte(uid.Value))
			uid.Value = hex.EncodeToString(checksum[:])
		}
	}
}
