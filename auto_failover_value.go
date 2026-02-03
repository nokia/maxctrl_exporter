// Copyright 2019, Vitaly Bezgachev, vitaly.bezgachev [the_at_symbol] gmail.com, Kadir Tugan, kadir.tugan [the_at_symbol] gmail.com
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"strings"
)

// AutoFailoverValue is a custom type that can handle both bool and string values
// String values are expected to be one of: "true", "false", or "safe"
type AutoFailoverValue struct {
	StringValue string
	BoolValue   bool
	IsString    bool
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *AutoFailoverValue) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as bool first
	var boolValue bool
	if err := json.Unmarshal(data, &boolValue); err == nil {
		a.BoolValue = boolValue
		a.StringValue = ""
		a.IsString = false
		return nil
	}

	// If not a bool, try as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		// Validate that the string is one of the expected enum values
		lowerValue := strings.ToLower(stringValue)
		if lowerValue == "true" || lowerValue == "false" || lowerValue == "safe" {
			a.StringValue = stringValue
			a.BoolValue = false
			a.IsString = true
			return nil
		} else {
			// If not a valid enum value, default to "false"
			a.StringValue = "false"
			a.BoolValue = false
			a.IsString = true
			return nil
		}
	}

	// If neither, return error
	return json.Unmarshal(data, nil)
}

// MarshalJSON implements the json.Marshaler interface
func (a AutoFailoverValue) MarshalJSON() ([]byte, error) {
	if a.IsString {
		return json.Marshal(a.StringValue)
	}
	return json.Marshal(a.BoolValue)
}

// String returns the string representation
// For boolean values: true becomes "true", false becomes "false"
// For string values: returns the original string value ("true", "false", or "safe")
func (a AutoFailoverValue) String() string {
	if a.IsString {
		return a.StringValue
	}
	if a.BoolValue {
		return "true"
	}
	return "false"
}

// Bool returns the boolean value
// For string values: "true" returns true, "safe" returns true, all others return false
func (a AutoFailoverValue) Bool() bool {
	if a.IsString {
		s := strings.ToLower(a.StringValue)
		return s == "true" || s == "safe"
	}
	return a.BoolValue
}
