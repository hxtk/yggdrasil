package urn

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Unstringer is an interface implemented by types that can parse themselves from strings.
type Unstringer interface {
	// FromString is a mutating operation on the receiver which converts it to
	// the value parsed from the string passed in.
	//
	// This operation should be idempotent and independent of the initial state
	// of the receiver. The receiver must be a pointer or this operation will
	// have no effect.
	//
	// If parsing fails, return an error shall be returned. Otherwise, the return
	// value shall be `nil`.
	FromString(string) error
}

// URN represents a Relative Resource Name as defined by Google's Resource-Oriented Design Guide.
//
// A URN should not have a leading or trailing `/`, and should consist only of characters
// which are valid within a URL.
type URN struct {
	// Parts represents the `/`-delimited components of the URN, with `/`s omitted.
	Parts []string
}

// Parse parses a URN into its components from a string.
func Parse(urn string) URN {
	return URN{Parts: strings.Split(urn, "/")}
}

// Scan parses the Parts of a URN into the destinations provided.
//
// Any number of destinations may be provided up to or including the
// number of segments in the URN.
//
// To consume values without reading them, `nil` receivers may be used
//
// Receivers must be pointers to `string` or `int`, or must implement
// `Unstringer`.
//
// An error will be returned if any non-nil receiver is not a pointer,
// if the number of receivers exceeds the number of parts, or if there
// is an error encountered while parsing the URN Part into the type of
// the corresponding receiver.
//
// Note that if an error is encountered during scanning, it should be
// assumed that any or all of the receivers may have been modified
// from their original values, but those values should not be considered
// valid.
func (u URN) Scan(dest ...interface{}) error {
	if len(dest) > len(u.Parts) {
		return errors.New("urn: fewer parts than receivers")
	}

	for i, v := range dest {
		if v == nil {
			continue
		}

		vType := reflect.TypeOf(v)
		if vType.Kind() != reflect.Ptr {
			return errors.New("urn: receivers must be pointer types")
		}

		switch x := v.(type) {
		case nil:
			continue
		case *string:
			*x = u.Parts[i]
		case *int:
			tmp, err := strconv.Atoi(u.Parts[i])
			if err != nil {
				return fmt.Errorf("urn: %v", err)
			}
			*x = tmp
		case *int64:
			tmp, err := strconv.ParseInt(u.Parts[i], 10, 0)
			if err != nil {
				return fmt.Errorf("urn: %v", err)
			}
			*x = tmp
		case *int32:
			tmp, err := strconv.ParseInt(u.Parts[i], 10, 0)
			if err != nil {
				return fmt.Errorf("urn: %v", err)
			}
			*x = int32(tmp)
		case Unstringer:
			err := x.FromString(u.Parts[i])
			if err != nil {
				return fmt.Errorf("urn: %v", err)
			}
		default:
			return errors.New("urn: unrecognized receiver type")
		}
	}

	return nil
}
