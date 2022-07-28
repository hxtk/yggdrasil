package authz

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestGetResourceName(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		msg  proto.Message
	}{
		{
			name: "foo",
			msg: &GetBookRequest{
				Name: "foo",
			},
		},
		{
			name: "foo",
			msg: &CreateBookRequest{
				Book: &Book{
					Name: "foo",
				},
			},
		},
		{
			name: "foo",
			msg: &ListBooksRequest{
				Parent: "foo",
			},
		},
		{
			name: "foo",
			msg: &UpdateBookRequest{
				Book: &Book{
					Name: "foo",
				},
			},
		},
	}

	for _, v := range testCases {
		name, err := getResourceName(v.msg.ProtoReflect())
		if err != v.err {
			t.Errorf("Expected error %v; got %v", v.err, err)
		}
		if name != v.name {
			t.Errorf("Expected name %v; got %v", v.name, name)
		}
	}
}
