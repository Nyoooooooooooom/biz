package output

import "testing"

func TestOKEnvelope(t *testing.T) {
	env := OK("trace-1", "done", map[string]string{"a": "b"})
	if env.APIVersion != "v1" || env.Code != "OK" || env.TraceID != "trace-1" {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}
