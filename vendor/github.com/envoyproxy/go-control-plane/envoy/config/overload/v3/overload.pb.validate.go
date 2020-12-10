// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/config/overload/v3/overload.proto

package envoy_config_overload_v3

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

	"github.com/golang/protobuf/ptypes"
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
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _overload_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on ResourceMonitor with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ResourceMonitor) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		return ResourceMonitorValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
	}

	switch m.ConfigType.(type) {

	case *ResourceMonitor_TypedConfig:

		if v, ok := interface{}(m.GetTypedConfig()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ResourceMonitorValidationError{
					field:  "TypedConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *ResourceMonitor_HiddenEnvoyDeprecatedConfig:

		if v, ok := interface{}(m.GetHiddenEnvoyDeprecatedConfig()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ResourceMonitorValidationError{
					field:  "HiddenEnvoyDeprecatedConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ResourceMonitorValidationError is the validation error returned by
// ResourceMonitor.Validate if the designated constraints aren't met.
type ResourceMonitorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ResourceMonitorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ResourceMonitorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ResourceMonitorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ResourceMonitorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ResourceMonitorValidationError) ErrorName() string { return "ResourceMonitorValidationError" }

// Error satisfies the builtin error interface
func (e ResourceMonitorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sResourceMonitor.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ResourceMonitorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ResourceMonitorValidationError{}

// Validate checks the field values on ThresholdTrigger with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ThresholdTrigger) Validate() error {
	if m == nil {
		return nil
	}

	if val := m.GetValue(); val < 0 || val > 1 {
		return ThresholdTriggerValidationError{
			field:  "Value",
			reason: "value must be inside range [0, 1]",
		}
	}

	return nil
}

// ThresholdTriggerValidationError is the validation error returned by
// ThresholdTrigger.Validate if the designated constraints aren't met.
type ThresholdTriggerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ThresholdTriggerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ThresholdTriggerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ThresholdTriggerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ThresholdTriggerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ThresholdTriggerValidationError) ErrorName() string { return "ThresholdTriggerValidationError" }

// Error satisfies the builtin error interface
func (e ThresholdTriggerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sThresholdTrigger.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ThresholdTriggerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ThresholdTriggerValidationError{}

// Validate checks the field values on ScaledTrigger with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ScaledTrigger) Validate() error {
	if m == nil {
		return nil
	}

	if val := m.GetScalingThreshold(); val < 0 || val > 1 {
		return ScaledTriggerValidationError{
			field:  "ScalingThreshold",
			reason: "value must be inside range [0, 1]",
		}
	}

	if val := m.GetSaturationThreshold(); val < 0 || val > 1 {
		return ScaledTriggerValidationError{
			field:  "SaturationThreshold",
			reason: "value must be inside range [0, 1]",
		}
	}

	return nil
}

// ScaledTriggerValidationError is the validation error returned by
// ScaledTrigger.Validate if the designated constraints aren't met.
type ScaledTriggerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ScaledTriggerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ScaledTriggerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ScaledTriggerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ScaledTriggerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ScaledTriggerValidationError) ErrorName() string { return "ScaledTriggerValidationError" }

// Error satisfies the builtin error interface
func (e ScaledTriggerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sScaledTrigger.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ScaledTriggerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ScaledTriggerValidationError{}

// Validate checks the field values on Trigger with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Trigger) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		return TriggerValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
	}

	switch m.TriggerOneof.(type) {

	case *Trigger_Threshold:

		if v, ok := interface{}(m.GetThreshold()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TriggerValidationError{
					field:  "Threshold",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Trigger_Scaled:

		if v, ok := interface{}(m.GetScaled()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TriggerValidationError{
					field:  "Scaled",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		return TriggerValidationError{
			field:  "TriggerOneof",
			reason: "value is required",
		}

	}

	return nil
}

// TriggerValidationError is the validation error returned by Trigger.Validate
// if the designated constraints aren't met.
type TriggerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TriggerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TriggerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TriggerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TriggerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TriggerValidationError) ErrorName() string { return "TriggerValidationError" }

// Error satisfies the builtin error interface
func (e TriggerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTrigger.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TriggerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TriggerValidationError{}

// Validate checks the field values on ScaleTimersOverloadActionConfig with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ScaleTimersOverloadActionConfig) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetTimerScaleFactors()) < 1 {
		return ScaleTimersOverloadActionConfigValidationError{
			field:  "TimerScaleFactors",
			reason: "value must contain at least 1 item(s)",
		}
	}

	for idx, item := range m.GetTimerScaleFactors() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ScaleTimersOverloadActionConfigValidationError{
					field:  fmt.Sprintf("TimerScaleFactors[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ScaleTimersOverloadActionConfigValidationError is the validation error
// returned by ScaleTimersOverloadActionConfig.Validate if the designated
// constraints aren't met.
type ScaleTimersOverloadActionConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ScaleTimersOverloadActionConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ScaleTimersOverloadActionConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ScaleTimersOverloadActionConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ScaleTimersOverloadActionConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ScaleTimersOverloadActionConfigValidationError) ErrorName() string {
	return "ScaleTimersOverloadActionConfigValidationError"
}

// Error satisfies the builtin error interface
func (e ScaleTimersOverloadActionConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sScaleTimersOverloadActionConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ScaleTimersOverloadActionConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ScaleTimersOverloadActionConfigValidationError{}

// Validate checks the field values on OverloadAction with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *OverloadAction) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		return OverloadActionValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
	}

	if len(m.GetTriggers()) < 1 {
		return OverloadActionValidationError{
			field:  "Triggers",
			reason: "value must contain at least 1 item(s)",
		}
	}

	for idx, item := range m.GetTriggers() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OverloadActionValidationError{
					field:  fmt.Sprintf("Triggers[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if v, ok := interface{}(m.GetTypedConfig()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return OverloadActionValidationError{
				field:  "TypedConfig",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// OverloadActionValidationError is the validation error returned by
// OverloadAction.Validate if the designated constraints aren't met.
type OverloadActionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OverloadActionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OverloadActionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OverloadActionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OverloadActionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OverloadActionValidationError) ErrorName() string { return "OverloadActionValidationError" }

// Error satisfies the builtin error interface
func (e OverloadActionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOverloadAction.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OverloadActionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OverloadActionValidationError{}

// Validate checks the field values on OverloadManager with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *OverloadManager) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetRefreshInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return OverloadManagerValidationError{
				field:  "RefreshInterval",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(m.GetResourceMonitors()) < 1 {
		return OverloadManagerValidationError{
			field:  "ResourceMonitors",
			reason: "value must contain at least 1 item(s)",
		}
	}

	for idx, item := range m.GetResourceMonitors() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OverloadManagerValidationError{
					field:  fmt.Sprintf("ResourceMonitors[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetActions() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OverloadManagerValidationError{
					field:  fmt.Sprintf("Actions[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// OverloadManagerValidationError is the validation error returned by
// OverloadManager.Validate if the designated constraints aren't met.
type OverloadManagerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OverloadManagerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OverloadManagerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OverloadManagerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OverloadManagerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OverloadManagerValidationError) ErrorName() string { return "OverloadManagerValidationError" }

// Error satisfies the builtin error interface
func (e OverloadManagerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOverloadManager.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OverloadManagerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OverloadManagerValidationError{}

// Validate checks the field values on
// ScaleTimersOverloadActionConfig_ScaleTimer with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *ScaleTimersOverloadActionConfig_ScaleTimer) Validate() error {
	if m == nil {
		return nil
	}

	if _, ok := _ScaleTimersOverloadActionConfig_ScaleTimer_Timer_NotInLookup[m.GetTimer()]; ok {
		return ScaleTimersOverloadActionConfig_ScaleTimerValidationError{
			field:  "Timer",
			reason: "value must not be in list [0]",
		}
	}

	if _, ok := ScaleTimersOverloadActionConfig_TimerType_name[int32(m.GetTimer())]; !ok {
		return ScaleTimersOverloadActionConfig_ScaleTimerValidationError{
			field:  "Timer",
			reason: "value must be one of the defined enum values",
		}
	}

	switch m.OverloadAdjust.(type) {

	case *ScaleTimersOverloadActionConfig_ScaleTimer_MinTimeout:

		if v, ok := interface{}(m.GetMinTimeout()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ScaleTimersOverloadActionConfig_ScaleTimerValidationError{
					field:  "MinTimeout",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *ScaleTimersOverloadActionConfig_ScaleTimer_MinScale:

		if v, ok := interface{}(m.GetMinScale()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ScaleTimersOverloadActionConfig_ScaleTimerValidationError{
					field:  "MinScale",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		return ScaleTimersOverloadActionConfig_ScaleTimerValidationError{
			field:  "OverloadAdjust",
			reason: "value is required",
		}

	}

	return nil
}

// ScaleTimersOverloadActionConfig_ScaleTimerValidationError is the validation
// error returned by ScaleTimersOverloadActionConfig_ScaleTimer.Validate if
// the designated constraints aren't met.
type ScaleTimersOverloadActionConfig_ScaleTimerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ScaleTimersOverloadActionConfig_ScaleTimerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ScaleTimersOverloadActionConfig_ScaleTimerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ScaleTimersOverloadActionConfig_ScaleTimerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ScaleTimersOverloadActionConfig_ScaleTimerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ScaleTimersOverloadActionConfig_ScaleTimerValidationError) ErrorName() string {
	return "ScaleTimersOverloadActionConfig_ScaleTimerValidationError"
}

// Error satisfies the builtin error interface
func (e ScaleTimersOverloadActionConfig_ScaleTimerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sScaleTimersOverloadActionConfig_ScaleTimer.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ScaleTimersOverloadActionConfig_ScaleTimerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ScaleTimersOverloadActionConfig_ScaleTimerValidationError{}

var _ScaleTimersOverloadActionConfig_ScaleTimer_Timer_NotInLookup = map[ScaleTimersOverloadActionConfig_TimerType]struct{}{
	0: {},
}
