package reasonerrors

import (
	"errors"
	"fmt"
	"strconv"

	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReasonNo is unknown reason number for error info.
	UnknownReasonNo = 0
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

//go:generate protoc -I. --go_out=paths=source_relative:. errors.proto

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason_no=%d reason = %s message = %s metadata = %v", e.Code, e.ReasonNo, e.Reason, e.Message, e.Metadata)
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata["reason_no"] = strconv.Itoa(int(e.ReasonNo))

	s, _ := status.New(httpstatus.ToGRPCCode(int(e.Code)), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Reason == e.Reason
	}
	return false
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := proto.Clone(e).(*Error)
	err.Metadata = md
	return err
}

// New returns an error object for the code, message.
func New(code, reasonNo int, reason, message string) *Error {
	if message == "" {
		message = reason
	}

	return &Error{
		Code:     int32(code),
		Message:  message,
		ReasonNo: int32(reasonNo),
		Reason:   reason,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code, reasonNo int, reason, format string, a ...interface{}) *Error {
	return New(code, reasonNo, reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code, reasonNo int, reason, format string, a ...interface{}) error {
	return New(code, reasonNo, reason, fmt.Sprintf(format, a...))
}

// Code returns the http code for a error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200 //nolint:gomnd
	}
	return int(FromError(err).Code)
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// ReasonNo returns the reason no for a particular error.
// It supports wrapped errors.
func ReasonNo(err error) int {
	if err == nil {
		return UnknownReasonNo
	}

	return int(FromError(err).ReasonNo)
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		ret := New(
			httpstatus.FromGRPCCode(gs.Code()),
			UnknownReasonNo,
			UnknownReason,
			gs.Message(),
		)
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				ret.Reason = d.Reason

				if v, ok := d.Metadata["reason_no"]; ok {
					if i, err := strconv.Atoi(v); err == nil {
						ret.ReasonNo = int32(i)
					}
				}

				return ret.WithMetadata(d.Metadata)
			}
		}
		return ret
	}
	return New(UnknownCode, UnknownReasonNo, UnknownReason, err.Error())
}
