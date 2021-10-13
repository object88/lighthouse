package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigSpec is the spec for a Config
type ConfigSpec struct {
	ReleaseName string `json:"releasename"`

	// Match is a list of match items which consist of select queries and expected match values or regular expressions.
	// When all match items for an object are positive, the rule is in effect.
	// +kubebuilder:validation:MinItems=1
	Match []MatchItem `json:"match"`
}

// MatchItem represents a single match query.
type MatchItem struct {
	// Select is a JSONPath query expression: https://goessner.net/articles/JsonPath/ which yields zero or more values.
	// If no match value or regex is specified, if the query yields a non-empty result, the match is considered positive.
	Select string `json:"select"`

	// MatchFor instructs how to match the results against the match... requirements.
	// Valid values are:
	// - "Any" - the match is considered positive if any of the results of select have a match.
	// - "All" - the match is considered positive only if all of the results of select have a match.
	// +optional
	MatchFor MatchForType `json:"matchFor,omitempty"`

	// MatchValue specifies the exact value to match the result of Select by.
	// The match is considered positive if at least one of the results of evaluating the select query yields a match when compared to matchValue.
	// +nullable
	MatchValue *string `json:"matchValue,omitempty"`

	// MatchValues specifies a list of values to match the result of Select by.
	// The match is considered positive if at least one of the results of evaluating the select query yields a match when compared to any of the values in the array.
	// +optional
	MatchValues []string `json:"matchValues,omitempty"`

	// MatchRegex specifies the regular expression to compare the result of Select by.
	// The match is considered positive if at least one of the results of evaluating the select query yields a match when compared to value.
	// +nullable
	MatchRegex *string `json:"matchRegex,omitempty"`

	// Negate indicates whether the match result should be to inverted.
	// Defaults to false.
	// +optional
	Negate bool `json:"negate,omitempty"`
}

// ModRuleType describes the type of a ModRule.
// Only one of the following ModRule types may be specified.
// +kubebuilder:validation:Enum=Patch;Reject
type ModRuleType string

const (
	// ModRuleTypePatch describes a ModRule which performs modifications on the target resource.
	ModRuleTypePatch ModRuleType = "Patch"

	// ModRuleTypeReject indicates that the ModRule should reject Create events for resources which match the rule.
	ModRuleTypeReject ModRuleType = "Reject"
)

// MatchForType describes the type of a match.
// Only one of the following ModRule types may be specified.
// +kubebuilder:validation:Enum=Any;All
type MatchForType string

const (
	// MatchForTypeAny indicates that a match is positive when any of the selected results matches any of the match requirements.
	MatchForTypeAny MatchForType = "Any"
	// MatchForTypeAll indicates that a match is positive when all of the selected results matches any of the match requirements.
	MatchForTypeAll MatchForType = "All"
)

// ConfigStatus is the status for a Config resource
type ConfigStatus struct {
	DeployedAt metav1.Time `json:"deployedat"`
}

// Config describes a target configuration
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigSpec   `json:"spec,omitempty"`
	Status ConfigStatus `json:"status,omitempty"`
}

// ConfigList is a list of Config resources
//+kubebuilder:object:root=true
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Config `json:"items"`
}
