package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"sigs.k8s.io/yaml"

	"github.com/authgear/authgear-server/pkg/util/validation"
)

type ProviderType string

const (
	ProviderTypeTwilio    ProviderType = "twilio"
	ProviderTypeNexmo     ProviderType = "nexmo"
	ProviderTypeAccessYou ProviderType = "accessyou"
	ProviderTypeSendCloud ProviderType = "sendcloud"
)

var _ = SMSProviderConfigSchema.Add("ProviderType", `
{
	"type": "string",
	"enum": ["twilio", "nexmo", "accessyou", "sendcloud"]
}
`)

type Provider struct {
	Name      string                   `json:"name,omitempty"`
	Type      ProviderType             `json:"type,omitempty"`
	Twilio    *ProviderConfigTwilio    `json:"twilio,omitempty" nullable:"true"`
	Nexmo     *ProviderConfigNexmo     `json:"nexmo,omitempty" nullable:"true"`
	AccessYou *ProviderConfigAccessYou `json:"accessyou,omitempty" nullable:"true"`
	SendCloud *ProviderConfigSendCloud `json:"sendcloud,omitempty" nullable:"true"`
}

type ProviderConfigTwilio struct {
	Sender              string `json:"sender,omitempty"`
	AccountSID          string `json:"account_sid,omitempty"`
	AuthToken           string `json:"auth_token,omitempty"`
	MessagingServiceSID string `json:"message_service_sid,omitempty"`
}

type ProviderConfigNexmo struct {
	Sender    string `json:"sender,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
	APISecret string `json:"api_secret,omitempty"`
}

type ProviderConfigAccessYou struct {
	Sender    string `json:"sender,omitempty"`
	BaseUrl   string `json:"base_url,omitempty"`
	AccountNo string `json:"accountno,omitempty"`
	User      string `json:"user,omitempty"`
	Pwd       string `json:"pwd,omitempty"`
}

var _ = SMSProviderConfigSchema.Add("Provider", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"name": { "type": "string" },
		"type": { "$ref": "#/$defs/ProviderType" },
		"twilio": { "$ref": "#/$defs/ProviderConfigTwilio" },
		"nexmo": { "$ref": "#/$defs/ProviderConfigNexmo" },
		"accessyou": { "$ref": "#/$defs/ProviderConfigAccessYou" },
		"sendcloud": { "$ref": "#/$defs/ProviderConfigSendCloud" }
	},
	"allOf": [
		{
			"if": { "properties": { "type": { "const": "twilio" } }},
			"then": { "required": ["twilio"] }
		},
		{
			"if": { "properties": { "type": { "const": "nexmo" } }},
			"then": { "required": ["nexmo"] }
		},
		{
			"if": { "properties": { "type": { "const": "accessyou" } }},
			"then": { "required": ["accessyou"] }
		},
		{
			"if": { "properties": { "type": { "const": "sendcloud" } }},
			"then": { "required": ["sendcloud"] }
		}
	]
}
`)

var _ = SMSProviderConfigSchema.Add("ProviderConfigTwilio", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"sender": { "type": "string" },
		"account_sid": { "type": "string" },
		"auth_token": {"type": "string"},
		"message_service_sid": {"type": "string"}
	},
	"required": ["sender", "account_sid", "auth_token", "message_service_sid"]
}
`)

var _ = SMSProviderConfigSchema.Add("ProviderConfigNexmo", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"sender": { "type": "string" },
		"api_key": { "type": "string" },
		"api_secret": {"type": "string"}
	},
	"required": ["sender", "api_key", "api_secret"]
}
`)

var _ = SMSProviderConfigSchema.Add("ProviderConfigAccessYou", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"sender": { "type": "string" },
		"base_url": { "type": "string" },
		"accountno": { "type": "string" },
		"user": { "type": "string" },
		"pwd": {"type": "string"}
	},
	"required": ["sender", "accountno", "user", "pwd"]
}
`)

type ProviderSelectorSwitchType string

const (
	ProviderSelectorSwitchTypeMatchPhoneNumberAlpha2         ProviderSelectorSwitchType = "match_phone_number_alpha2"
	ProviderSelectorSwitchTypeMatchAppIDAndPhoneNumberAlpha2 ProviderSelectorSwitchType = "match_app_id_and_phone_number_alpha2"
	ProviderSelectorSwitchTypeDefault                        ProviderSelectorSwitchType = "default"
)

var _ = SMSProviderConfigSchema.Add("ProviderSelectorSwitchType", `
{
	"type": "string",
	"enum": ["match_phone_number_alpha2", "match_app_id_and_phone_number_alpha2", "default"]
}
`)

type ProviderSelectorSwitchRule struct {
	Type              ProviderSelectorSwitchType `json:"type,omitempty"`
	UseProvider       string                     `json:"use_provider,omitempty"`
	PhoneNumberAlpha2 string                     `json:"phone_number_alpha2,omitempty"`
	AppID             string                     `json:"app_id,omitempty"`
}

var _ = SMSProviderConfigSchema.Add("ProviderSelectorSwitchRule", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"type": { "$ref": "#/$defs/ProviderSelectorSwitchType" },
		"use_provider": { "type": "string" },
		"phone_number_alpha2": { "type": "string" },
		"app_id": { "type": "string" }
	},
	"allOf": [
		{
			"if": { "properties": { "type": { "const": "phone_number_alpha2" } }},
			"then": { "required": ["phone_number_alpha2"] }
		},
		{
			"if": { "properties": { "type": { "const": "app_id_and_phone_number_alpha2" } }},
			"then": { "required": ["phone_number_alpha2"] }
		}
	],
	"required": ["type", "use_provider"]
}
`)

type ProviderSelector struct {
	Switch []*ProviderSelectorSwitchRule `json:"switch,omitempty"`
}

var _ = SMSProviderConfigSchema.Add("ProviderSelector", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"switch": {
			"type": "array",
			"minItems": 1,
			"items": { "$ref": "#/$defs/ProviderSelectorSwitchRule" }
		}
	},
	"required": ["switch"]
}
`)

type SMSProviderConfig struct {
	Providers        []*Provider       `json:"providers,omitempty"`
	ProviderSelector *ProviderSelector `json:"provider_selector,omitempty"`
}

var _ = SMSProviderConfigSchema.Add("SMSProviderConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"providers": {
			"type": "array",
			"minItems": 1,
			"items": { "$ref": "#/$defs/Provider" }
		},
		"provider_selector": { "$ref": "#/$defs/ProviderSelector" }
	},
	"required": ["providers", "provider_selector"]
}
`)

func (c *SMSProviderConfig) Validate(ctx *validation.Context) {
	c.ValidateProviderSelectorUseProvider(ctx)
	c.ValidateProviderSelectorDefault(ctx)
	c.ValidateSendCloudConfigs(ctx)
}

func (c *SMSProviderConfig) ValidateProviderSelectorUseProvider(ctx *validation.Context) {
	providers := c.Providers
	for i, switchCase := range c.ProviderSelector.Switch {
		useProvider := switchCase.UseProvider
		idx := slices.IndexFunc(providers, func(p *Provider) bool { return p.Name == useProvider })
		if idx == -1 {
			ctx.Child("provider_selector", "switch", strconv.Itoa(i), "use_provider").EmitErrorMessage(fmt.Sprintf("provider %s not found", useProvider))
		}
	}
}

func (c *SMSProviderConfig) ValidateProviderSelectorDefault(ctx *validation.Context) {
	for _, switchCase := range c.ProviderSelector.Switch {
		if switchCase.Type == ProviderSelectorSwitchTypeDefault {
			return
		}
	}
	ctx.Child("provider_selector", "switch").EmitErrorMessage(fmt.Sprintf("provider selector default not found"))
}

func (c *SMSProviderConfig) ValidateSendCloudConfigs(ctx *validation.Context) {
	for i, provider := range c.Providers {
		if provider.Type == ProviderTypeSendCloud {
			c.ValidateSendCloudConfig(ctx.Child("providers", strconv.Itoa(i), "sendcloud"), provider.SendCloud)
		}
	}
}

func (c *SMSProviderConfig) ValidateSendCloudConfig(ctx *validation.Context, sendCloudConfig *ProviderConfigSendCloud) {
	templates := sendCloudConfig.Templates
	for i, templateAssignment := range sendCloudConfig.TemplateAssignments {
		ctxTemplateAssignment := ctx.Child("template_assignments", strconv.Itoa(i))
		defaultTemplateID := templateAssignment.DefaultTemplateID
		idx := slices.IndexFunc(templates, func(t *SendCloudTemplate) bool { return t.TemplateID == defaultTemplateID })

		if idx == -1 {
			ctxTemplateAssignment.Child("default_template_id").EmitErrorMessage(fmt.Sprintf("template_id %v not found", defaultTemplateID))
		}

		for j, byLanguage := range templateAssignment.ByLanguages {
			ctxByLanguage := ctxTemplateAssignment.Child("by_languages", strconv.Itoa(j))
			templateID := byLanguage.TemplateID
			idx = slices.IndexFunc(templates, func(t *SendCloudTemplate) bool { return t.TemplateID == templateID })
			if idx == -1 {
				ctxByLanguage.Child("template_id").EmitErrorMessage(fmt.Sprintf("template_id %v not found", templateID))
			}
		}
	}

}

func ParseSMSProviderConfigFromYAML(inputYAML []byte) (*SMSProviderConfig, error) {
	const validationErrorMessage = "invalid configuration"

	jsonData, err := yaml.YAMLToJSON(inputYAML)
	if err != nil {
		return nil, err
	}

	err = SMSProviderConfigSchema.Validator().ValidateWithMessage(bytes.NewReader(jsonData), validationErrorMessage)
	if err != nil {
		return nil, err
	}

	var config SMSProviderConfig
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	err = validation.ValidateValueWithMessage(&config, validationErrorMessage)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
