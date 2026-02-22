package routes

type ApiError struct {
	Code       string           `json:"code,omitempty"`
	Target     *string          `json:"target,omitempty"`
	Message    string           `json:"message,omitempty"`
	InnerError *ApiError        `json:"innerError,omitempty"`
	Details    []ApiErrorDetail `json:"details,omitempty"`
}

type ApiErrorDetail struct {
	Code    *string `json:"code,omitempty"`
	Message string  `json:"message,omitempty"`
	Target  string  `json:"target,omitempty"`
}

type Response struct {
	Ok       bool      `json:"ok"`
	Error    *ApiError `json:"error,omitempty"`
	TraceId  *string   `json:"traceId,omitempty"`
	SpanId   *string   `json:"spanId,omitempty"`
	Value    any       `json:"value,omitempty"`
	NextLink *string   `json:"@nextLink,omitempty"`
}

func SuccessResponse(value any) Response {
	return Response{
		Ok:    true,
		Value: value,
	}
}

func ErrorResponse(err *ApiError) Response {
	return Response{
		Ok:    false,
		Error: err,
	}
}
