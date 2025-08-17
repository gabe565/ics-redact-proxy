package ics

import (
	"crypto/sha256"
	"encoding/hex"
	"slices"

	"gabe565.com/ics-redact-proxy/internal/config"
	ics "github.com/arran4/golang-ical"
)

func FilterComponent(conf *config.Config, c *ics.ComponentBase) {
	c.Properties = slices.DeleteFunc(c.Properties, func(property ics.IANAProperty) bool {
		return !slices.Contains(conf.ComponentFields, property.IANAToken)
	})
	c.Properties = slices.Clip(c.Properties)

	if conf.NewEventSummary != "" {
		c.SetSummary(conf.NewEventSummary)
	}

	if conf.HashUID {
		if uid := c.GetProperty(ics.ComponentPropertyUniqueId); uid != nil {
			checksum := sha256.Sum256([]byte(uid.Value))
			uid.Value = hex.EncodeToString(checksum[:])
		}
	}
}
