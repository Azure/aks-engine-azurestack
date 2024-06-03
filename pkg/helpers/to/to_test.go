// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package to

import "testing"

func TestString(t *testing.T) {
	v := ""
	if String(&v) != v {
		t.Fatalf("to: String failed to return the correct string -- expected %v, received %v",
			v, String(&v))
	}
}

func TestStringHandlesNil(t *testing.T) {
	if String(nil) != "" {
		t.Fatalf("to: String failed to correctly convert nil -- expected %v, received %v",
			"", String(nil))
	}
}

func TestStringPtr(t *testing.T) {
	v := ""
	if *StringPtr(v) != v {
		t.Fatalf("to: StringPtr failed to return the correct string -- expected %v, received %v",
			v, *StringPtr(v))
	}
}

func TestBool(t *testing.T) {
	v := false
	if Bool(&v) != v {
		t.Fatalf("to: Bool failed to return the correct string -- expected %v, received %v",
			v, Bool(&v))
	}
}

func TestBoolHandlesNil(t *testing.T) {
	if Bool(nil) != false {
		t.Fatalf("to: Bool failed to correctly convert nil -- expected %v, received %v",
			false, Bool(nil))
	}
}

func TestBoolPtr(t *testing.T) {
	v := false
	if *BoolPtr(v) != v {
		t.Fatalf("to: BoolPtr failed to return the correct string -- expected %v, received %v",
			v, *BoolPtr(v))
	}
}

func TestInt(t *testing.T) {
	v := 0
	if Int(&v) != v {
		t.Fatalf("to: Int failed to return the correct string -- expected %v, received %v",
			v, Int(&v))
	}
}

func TestIntHandlesNil(t *testing.T) {
	if Int(nil) != 0 {
		t.Fatalf("to: Int failed to correctly convert nil -- expected %v, received %v",
			0, Int(nil))
	}
}

func TestIntPtr(t *testing.T) {
	v := 0
	if *IntPtr(v) != v {
		t.Fatalf("to: IntPtr failed to return the correct string -- expected %v, received %v",
			v, *IntPtr(v))
	}
}

func TestInt32(t *testing.T) {
	v := int32(0)
	if Int32(&v) != v {
		t.Fatalf("to: Int32 failed to return the correct string -- expected %v, received %v",
			v, Int32(&v))
	}
}

func TestInt32HandlesNil(t *testing.T) {
	if Int32(nil) != int32(0) {
		t.Fatalf("to: Int32 failed to correctly convert nil -- expected %v, received %v",
			0, Int32(nil))
	}
}

func TestInt32Ptr(t *testing.T) {
	v := int32(0)
	if *Int32Ptr(v) != v {
		t.Fatalf("to: Int32Ptr failed to return the correct string -- expected %v, received %v",
			v, *Int32Ptr(v))
	}
}

func TestInt64(t *testing.T) {
	v := int64(0)
	if Int64(&v) != v {
		t.Fatalf("to: Int64 failed to return the correct string -- expected %v, received %v",
			v, Int64(&v))
	}
}

func TestInt64HandlesNil(t *testing.T) {
	if Int64(nil) != int64(0) {
		t.Fatalf("to: Int64 failed to correctly convert nil -- expected %v, received %v",
			0, Int64(nil))
	}
}

func TestInt64Ptr(t *testing.T) {
	v := int64(0)
	if *Int64Ptr(v) != v {
		t.Fatalf("to: Int64Ptr failed to return the correct string -- expected %v, received %v",
			v, *Int64Ptr(v))
	}
}

func TestFloat64(t *testing.T) {
	v := float64(0)
	if Float64(&v) != v {
		t.Fatalf("to: Float64 failed to return the correct string -- expected %v, received %v",
			v, Float64(&v))
	}
}

func TestFloat64HandlesNil(t *testing.T) {
	if Float64(nil) != float64(0) {
		t.Fatalf("to: Float64 failed to correctly convert nil -- expected %v, received %v",
			0, Float64(nil))
	}
}

func TestFloat64Ptr(t *testing.T) {
	v := float64(0)
	if *Float64Ptr(v) != v {
		t.Fatalf("to: Float64Ptr failed to return the correct string -- expected %v, received %v",
			v, *Float64Ptr(v))
	}
}
