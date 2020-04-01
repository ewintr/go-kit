package herror_test

import (
	"runtime"
	"testing"

	"git.sr.ht/~ewintr/go-kit/herror"
	"git.sr.ht/~ewintr/go-kit/test"
)

func trace() *herror.Stacktrace {
	return herror.NewStacktrace()
}

func traceStepIn(f []herror.FrameFilter) *herror.Stacktrace {
	return traceWithFilter(f)
}

func traceWithFilter(f []herror.FrameFilter) *herror.Stacktrace {
	return herror.NewStacktrace(f...)
}

func TestStacktrace(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		stack := trace()

		expectedFrames := []herror.Frame{
			herror.Frame{
				Function: "TestStacktrace.func1",
			},
			herror.Frame{
				Function: "trace",
			},
		}

		test.Equals(t, len(expectedFrames), len(stack.Frames))
		for i, frame := range expectedFrames {
			test.Equals(t, frame.Function, stack.Frames[i].Function)
			test.Equals(t, "git.sr.ht/~ewintr/go-kit/herror_test", stack.Frames[i].Package)
			test.Equals(t, "stacktrace_test.go", stack.Frames[i].Filename)
		}
	})

	t.Run("filter frames", func(t *testing.T) {

		for _, tc := range []struct {
			m        string
			filters  []herror.FrameFilter
			expected []herror.Frame
		}{
			{
				m: "no filter",
				expected: []herror.Frame{
					herror.Frame{
						Function: "TestStacktrace.func2",
					},
					herror.Frame{
						Function: "traceStepIn",
					},
					herror.Frame{
						Function: "traceWithFilter",
					},
				},
			},
			{
				m: "single filter",
				expected: []herror.Frame{
					herror.Frame{
						Function: "traceStepIn",
					},
					herror.Frame{
						Function: "traceWithFilter",
					},
				},
				filters: []herror.FrameFilter{
					func(f herror.Frame) bool {
						return f.Function == "TestStacktrace.func2"
					},
				},
			},
			{
				m: "multiple filters",
				expected: []herror.Frame{
					herror.Frame{
						Function: "traceWithFilter",
					},
				},
				filters: []herror.FrameFilter{
					func(f herror.Frame) bool {
						return f.Function == "TestStacktrace.func2"
					},
					func(f herror.Frame) bool {
						return f.Function == "traceStepIn"
					},
				},
			},
		} {
			stack := traceStepIn(tc.filters)

			t.Run(tc.m, func(t *testing.T) {
				test.Equals(t, len(tc.expected), len(stack.Frames))

				for i, frame := range tc.expected {
					test.Equals(t, frame.Function, stack.Frames[i].Function)
					test.Equals(t, "git.sr.ht/~ewintr/go-kit/herror_test", stack.Frames[i].Package)
					test.Equals(t, "stacktrace_test.go", stack.Frames[i].Filename)
				}
			})
		}
	})
}

func TestFrame(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		f := func() herror.Frame {
			pc := make([]uintptr, 1)
			n := runtime.Callers(0, pc)
			test.Assert(t, n == 1, "expected available pcs")

			frames := runtime.CallersFrames(pc)
			runtimeframe, _ := frames.Next()
			return herror.NewFrame(runtimeframe)
		}

		frame := f()
		test.Equals(t, "Callers", frame.Function)
		test.Equals(t, "runtime", frame.Package)
		test.Equals(t, "extern.go", frame.Filename)
	})
}
