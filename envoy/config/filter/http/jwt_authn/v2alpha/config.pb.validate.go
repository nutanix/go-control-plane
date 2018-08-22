// Code generated by protoc-gen-validate
// source: envoy/config/filter/http/jwt_authn/v2alpha/config.proto
// DO NOT EDIT!!!

package envoy_config_filter_http_jwt_authn_v2alpha

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gogo/protobuf/types"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = types.DynamicAny{}
)

// Validate checks the field values on JwtProvider with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *JwtProvider) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetIssuer()) < 1 {
		return JwtProviderValidationError{
			Field:  "Issuer",
			Reason: "value length must be at least 1 bytes",
		}
	}

	// no validation rules for Forward

	for idx, item := range m.GetFromHeaders() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtProviderValidationError{
					Field:  fmt.Sprintf("FromHeaders[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	// no validation rules for ForwardPayloadHeader

	switch m.JwksSourceSpecifier.(type) {

	case *JwtProvider_RemoteJwks:

		if v, ok := interface{}(m.GetRemoteJwks()).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtProviderValidationError{
					Field:  "RemoteJwks",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *JwtProvider_LocalJwks:

		if v, ok := interface{}(m.GetLocalJwks()).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtProviderValidationError{
					Field:  "LocalJwks",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	default:
		return JwtProviderValidationError{
			Field:  "JwksSourceSpecifier",
			Reason: "value is required",
		}

	}

	return nil
}

// JwtProviderValidationError is the validation error returned by
// JwtProvider.Validate if the designated constraints aren't met.
type JwtProviderValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e JwtProviderValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJwtProvider.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = JwtProviderValidationError{}

// Validate checks the field values on RemoteJwks with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *RemoteJwks) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetHttpUri()).(interface {
		Validate() error
	}); ok {
		if err := v.Validate(); err != nil {
			return RemoteJwksValidationError{
				Field:  "HttpUri",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetCacheDuration()).(interface {
		Validate() error
	}); ok {
		if err := v.Validate(); err != nil {
			return RemoteJwksValidationError{
				Field:  "CacheDuration",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// RemoteJwksValidationError is the validation error returned by
// RemoteJwks.Validate if the designated constraints aren't met.
type RemoteJwksValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e RemoteJwksValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoteJwks.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = RemoteJwksValidationError{}

// Validate checks the field values on JwtHeader with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *JwtHeader) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetName()) < 1 {
		return JwtHeaderValidationError{
			Field:  "Name",
			Reason: "value length must be at least 1 bytes",
		}
	}

	// no validation rules for ValuePrefix

	return nil
}

// JwtHeaderValidationError is the validation error returned by
// JwtHeader.Validate if the designated constraints aren't met.
type JwtHeaderValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e JwtHeaderValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJwtHeader.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = JwtHeaderValidationError{}

// Validate checks the field values on ProviderWithAudiences with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ProviderWithAudiences) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ProviderName

	return nil
}

// ProviderWithAudiencesValidationError is the validation error returned by
// ProviderWithAudiences.Validate if the designated constraints aren't met.
type ProviderWithAudiencesValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e ProviderWithAudiencesValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sProviderWithAudiences.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = ProviderWithAudiencesValidationError{}

// Validate checks the field values on JwtRequirement with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *JwtRequirement) Validate() error {
	if m == nil {
		return nil
	}

	switch m.RequiresType.(type) {

	case *JwtRequirement_ProviderName:
		// no validation rules for ProviderName

	case *JwtRequirement_ProviderAndAudiences:

		if v, ok := interface{}(m.GetProviderAndAudiences()).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtRequirementValidationError{
					Field:  "ProviderAndAudiences",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *JwtRequirement_RequiresAny:

		if v, ok := interface{}(m.GetRequiresAny()).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtRequirementValidationError{
					Field:  "RequiresAny",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *JwtRequirement_RequiresAll:

		if v, ok := interface{}(m.GetRequiresAll()).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtRequirementValidationError{
					Field:  "RequiresAll",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *JwtRequirement_AllowMissingOrFailed:

		if v, ok := interface{}(m.GetAllowMissingOrFailed()).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtRequirementValidationError{
					Field:  "AllowMissingOrFailed",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	return nil
}

// JwtRequirementValidationError is the validation error returned by
// JwtRequirement.Validate if the designated constraints aren't met.
type JwtRequirementValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e JwtRequirementValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJwtRequirement.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = JwtRequirementValidationError{}

// Validate checks the field values on JwtRequirementOrList with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *JwtRequirementOrList) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetRequirements()) < 2 {
		return JwtRequirementOrListValidationError{
			Field:  "Requirements",
			Reason: "value must contain at least 2 item(s)",
		}
	}

	for idx, item := range m.GetRequirements() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtRequirementOrListValidationError{
					Field:  fmt.Sprintf("Requirements[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	return nil
}

// JwtRequirementOrListValidationError is the validation error returned by
// JwtRequirementOrList.Validate if the designated constraints aren't met.
type JwtRequirementOrListValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e JwtRequirementOrListValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJwtRequirementOrList.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = JwtRequirementOrListValidationError{}

// Validate checks the field values on JwtRequirementAndList with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *JwtRequirementAndList) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetRequirements()) < 2 {
		return JwtRequirementAndListValidationError{
			Field:  "Requirements",
			Reason: "value must contain at least 2 item(s)",
		}
	}

	for idx, item := range m.GetRequirements() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtRequirementAndListValidationError{
					Field:  fmt.Sprintf("Requirements[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	return nil
}

// JwtRequirementAndListValidationError is the validation error returned by
// JwtRequirementAndList.Validate if the designated constraints aren't met.
type JwtRequirementAndListValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e JwtRequirementAndListValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJwtRequirementAndList.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = JwtRequirementAndListValidationError{}

// Validate checks the field values on RequirementRule with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *RequirementRule) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetMatch() == nil {
		return RequirementRuleValidationError{
			Field:  "Match",
			Reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetMatch()).(interface {
		Validate() error
	}); ok {
		if err := v.Validate(); err != nil {
			return RequirementRuleValidationError{
				Field:  "Match",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetRequires()).(interface {
		Validate() error
	}); ok {
		if err := v.Validate(); err != nil {
			return RequirementRuleValidationError{
				Field:  "Requires",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// RequirementRuleValidationError is the validation error returned by
// RequirementRule.Validate if the designated constraints aren't met.
type RequirementRuleValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e RequirementRuleValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRequirementRule.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = RequirementRuleValidationError{}

// Validate checks the field values on JwtAuthentication with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *JwtAuthentication) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Providers

	for idx, item := range m.GetRules() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface {
			Validate() error
		}); ok {
			if err := v.Validate(); err != nil {
				return JwtAuthenticationValidationError{
					Field:  fmt.Sprintf("Rules[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	return nil
}

// JwtAuthenticationValidationError is the validation error returned by
// JwtAuthentication.Validate if the designated constraints aren't met.
type JwtAuthenticationValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e JwtAuthenticationValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJwtAuthentication.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = JwtAuthenticationValidationError{}
