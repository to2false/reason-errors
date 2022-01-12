package errors

import (
	"testing"
)

func TestTypes(t *testing.T) {
	var (
		input = []error{
			BadRequest(0, "reason_400", "message_400"),
			Unauthorized(0, "reason_401", "message_401"),
			Forbidden(0, "reason_403", "message_403"),
			NotFound(0, "reason_404", "message_404"),
			Conflict(0, "reason_409", "message_409"),
			InternalServer(0, "reason_500", "message_500"),
			ServiceUnavailable(0, "reason_503", "message_503"),
			GatewayTimeout(0, "reason_504", "message_504"),
			ClientClosed(0, "reason_499", "message_499"),
		}
		output = []func(error) bool{
			IsBadRequest,
			IsUnauthorized,
			IsForbidden,
			IsNotFound,
			IsConflict,
			IsInternalServer,
			IsServiceUnavailable,
			IsGatewayTimeout,
			IsClientClosed,
		}
	)

	for i, in := range input {
		if !output[i](in) {
			t.Errorf("not expect: %v", in)
		}
	}
}
