package sms

import (
	"errors"
	"fmt"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

func GetClientNameByMatch(c *config.SMSProviderConfig, ctx *MatchContext) string {
	var defaultClient string
	for _, providerSelector := range c.ProviderSelector.Switch {
		matcher := ParseMatcher(providerSelector)
		switch m := matcher.(type) {
		case *MatcherDefault:
			defaultClient = providerSelector.UseProvider
			break
		default:
			if m.Match(ctx) {
				return providerSelector.UseProvider
			}
		}
	}
	if defaultClient == "" {
		panic(errors.New(fmt.Sprintf("Cannot select provider given %v", ctx)))
	}
	return defaultClient
}
