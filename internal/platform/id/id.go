package id

import "github.com/google/uuid"

func TraceID() string { return uuid.NewString() }
