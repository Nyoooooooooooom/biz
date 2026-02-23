package errors

import "testing"

func TestExitCode(t *testing.T) {
	cases := []struct {
		err  error
		want int
	}{
		{nil, 0},
		{New(KindValidation, "bad"), 2},
		{New(KindDependencyUnavailable, "down"), 3},
		{New(KindInternal, "oops"), 4},
	}
	for _, c := range cases {
		if got := ExitCode(c.err); got != c.want {
			t.Fatalf("got %d want %d", got, c.want)
		}
	}
}
