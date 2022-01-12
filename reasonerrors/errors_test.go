package reasonerrors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/grpc_testing"
)

func TestError(t *testing.T) {
	var base *Error
	err := Newf(http.StatusBadRequest, 1, "reason", "message")
	err2 := Newf(http.StatusBadRequest, 1, "reason", "message")
	err3 := err.WithMetadata(map[string]string{
		"foo": "bar",
	})
	werr := fmt.Errorf("wrap %w", err)

	if errors.Is(err, new(Error)) {
		t.Errorf("should not be equal: %v", err)
	}
	if !errors.Is(werr, err) {
		t.Errorf("should be equal: %v", err)
	}
	if !errors.Is(werr, err2) {
		t.Errorf("should be equal: %v", err)
	}

	if !errors.As(err, &base) {
		t.Errorf("should be matchs: %v", err)
	}
	if !IsBadRequest(err) {
		t.Errorf("should be matchs: %v", err)
	}

	if reason := Reason(err); reason != err3.Reason {
		t.Errorf("got %s want: %s", reason, err)
	}

	if err3.Metadata["foo"] != "bar" {
		t.Error("not expected metadata")
	}

	gs := err.GRPCStatus()
	se := FromError(gs.Err())
	if se.Reason != "reason" {
		t.Errorf("got %+v want %+v", se, err)
	}

	gs2, _ := status.New(codes.InvalidArgument, "bad request").WithDetails(&grpc_testing.Empty{})
	se2 := FromError(gs2.Err())
	// codes.InvalidArgument should convert to http.StatusBadRequest
	if se2.Code != http.StatusBadRequest {
		t.Errorf("convert code err, got %d want %d", UnknownCode, http.StatusBadRequest)
	}
	assert.Nil(t, FromError(nil))
	e := FromError(errors.New("test"))
	assert.Equal(t, e.Code, int32(UnknownCode))
}

func TestIs(t *testing.T) {
	tests := []struct {
		name string
		e    *Error
		err  error
		want bool
	}{
		{
			name: "true",
			e:    &Error{Reason: "test"},
			err:  New(http.StatusNotFound, 0, "test", ""),
			want: true,
		},
		{
			name: "false",
			e:    &Error{Reason: "test"},
			err:  errors.New("test"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.e.Is(tt.err); ok != tt.want {
				t.Errorf("Error.Error() = %v, want %v", ok, tt.want)
			}
		})
	}
}

func TestOther(t *testing.T) {
	assert.Equal(t, Code(nil), 200)
	assert.Equal(t, Code(errors.New("test")), UnknownCode)
	assert.Equal(t, Reason(errors.New("test")), UnknownReason)
	err := Errorf(10001, 0, "test code 10001", "message")
	assert.Equal(t, Code(err), 10001)
	assert.Equal(t, Reason(err), "test code 10001")
}
