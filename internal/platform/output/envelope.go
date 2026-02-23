package output

type Err struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

type Envelope[T any] struct {
	APIVersion string         `json:"api_version"`
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	TraceID    string         `json:"trace_id"`
	Data       T              `json:"data,omitempty"`
	Error      *Err           `json:"error,omitempty"`
	Meta       map[string]any `json:"meta,omitempty"`
}

func OK[T any](traceID, message string, data T) Envelope[T] {
	return Envelope[T]{
		APIVersion: "v1",
		Code:       "OK",
		Message:    message,
		TraceID:    traceID,
		Data:       data,
	}
}

func Fail[T any](traceID, code, message, kind string) Envelope[T] {
	return Envelope[T]{
		APIVersion: "v1",
		Code:       code,
		Message:    message,
		TraceID:    traceID,
		Error: &Err{
			Kind:    kind,
			Message: message,
		},
	}
}
