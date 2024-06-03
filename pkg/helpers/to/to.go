// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package to

// IsTrueBoolPointer is a simple boolean helper function for boolean pointers
func IsTrueBoolPointer(b *bool) bool {
	if b != nil && *b {
		return true
	}
	return false
}

// IsFalseBoolPointer is a simple boolean helper function for boolean pointers
func IsFalseBoolPointer(b *bool) bool {
	if b != nil && !*b {
		return true
	}
	return false
}

// Bool returns a bool value for the passed bool pointer. It returns false if the pointer is nil.
func Bool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

// String returns a string value for the passed string pointer. It returns the empty string if the
// pointer is nil.
func String(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// Int returns an int value for the passed int pointer. It returns 0 if the pointer is nil.
func Int(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

// Int32 returns an int value for the passed int pointer. It returns 0 if the pointer is nil.
func Int32(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

// Int64 returns an int value for the passed int pointer. It returns 0 if the pointer is nil.
func Int64(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

// Float64 returns an int value for the passed int pointer. It returns 0.0 if the pointer is nil.
func Float64(i *float64) float64 {
	if i != nil {
		return *i
	}
	return 0.0
}

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	p := b
	return &p
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	p := s
	return &p
}

// IntPtr returns a pointer to a int
func IntPtr(i int) *int {
	p := i
	return &p
}

// Int32Ptr returns a pointer to a int32
func Int32Ptr(i int32) *int32 {
	p := i
	return &p
}

// Int64Ptr returns a pointer to a int64
func Int64Ptr(i int64) *int64 {
	p := i
	return &p
}

// Float64Ptr returns a pointer to a float64
func Float64Ptr(i float64) *float64 {
	p := i
	return &p
}
