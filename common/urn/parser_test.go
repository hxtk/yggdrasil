package urn

import (
	"testing"
)

func TestScan(t *testing.T) {
	t.Run("parse int", func(t *testing.T) {
		res := Parse("1/2")
		var foo int
		var bar int
		err := res.Scan(&foo, &bar)
		if err != nil {
			t.Errorf("Scan encountered error: %v", err)
		}

		if foo != 1 {
			t.Errorf("Expected 1; got %v", foo)
		}

		if bar != 2 {
			t.Errorf("Expected 2; got %v", bar)
		}
	})
	t.Run("parse bad int", func(t *testing.T) {
		res := Parse("foo/2")
		var foo int
		var bar int
		err := res.Scan(&foo, &bar)
		if err == nil {
			t.Errorf("Scan did not error parsing \"foo\" as int")
		}
	})
	t.Run("parse int64", func(t *testing.T) {
		res := Parse("1/2")
		var foo int64
		var bar int64
		err := res.Scan(&foo, &bar)
		if err != nil {
			t.Errorf("Scan encountered error: %v", err)
		}

		if foo != 1 {
			t.Errorf("Expected 1; got %v", foo)
		}

		if bar != 2 {
			t.Errorf("Expected 2; got %v", bar)
		}
	})
	t.Run("parse int32", func(t *testing.T) {
		res := Parse("1/2")
		var foo int32
		var bar int32
		err := res.Scan(&foo, &bar)
		if err != nil {
			t.Errorf("Scan encountered error: %v", err)
		}

		if foo != 1 {
			t.Errorf("Expected 1; got %v", foo)
		}

		if bar != 2 {
			t.Errorf("Expected 2; got %v", bar)
		}
	})
	t.Run("parse string", func(t *testing.T) {
		res := Parse("foo/bar")
		var foo string
		var bar string
		err := res.Scan(&foo, &bar)
		if err != nil {
			t.Errorf("Scan encountered error: %v", err)
		}

		if foo != "foo" {
			t.Errorf("Expected \"foo\"; got %v", foo)
		}

		if bar != "bar" {
			t.Errorf("Expected \"bar\"; got %v", bar)
		}
	})
	t.Run("parse int and string", func(t *testing.T) {
		res := Parse("command/3")
		var foo string
		var bar int
		err := res.Scan(&foo, &bar)
		if err != nil {
			t.Errorf("Scan encountered error: %v", err)
		}

		if foo != "command" {
			t.Errorf("Expected \"command\"; got %v", foo)
		}

		if bar != 3 {
			t.Errorf("Expected 3; got %v", bar)
		}
	})
	t.Run("skip nil receivers", func(t *testing.T) {
		res := Parse("command/3")
		var bar int
		err := res.Scan(nil, &bar)
		if err != nil {
			t.Errorf("Scan encountered error: %v", err)
		}

		if bar != 3 {
			t.Errorf("Expected 3; got %v", bar)
		}
	})
	t.Run("too many receivers", func(t *testing.T) {
		res := Parse("command/3")
		var foo int
		var bar string
		err := res.Scan(nil, &foo, &bar)
		if err == nil {
			t.Errorf("Scan did not error parsing \"foo\" as int")
		}
	})
}
