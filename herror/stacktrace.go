package herror

import (
	"fmt"
	"go/build"
	"path/filepath"
	"runtime"
	"strings"
)

// Stacktrace holds information about the frames of the stack.
type Stacktrace struct {
	Frames []Frame `json:"frames,omitempty"`
}

// Frame represents parsed information from runtime.Frame
type Frame struct {
	Function string `json:"function,omitempty"`
	Type     string `json:"type,omitempty"`
	Package  string `json:"package,omitempty"`
	Filename string `json:"filename,omitempty"`
	AbsPath  string `json:"abs_path,omitempty"`
	Line     int    `json:"line,omitempty"`
	InApp    bool   `json:"in_app,omitempty"`
}

// FrameFilter represents function to filter frames
type FrameFilter func(Frame) bool

const unknown string = "unknown"

// NewStacktrace creates a stacktrace using `runtime.Callers`.
func NewStacktrace(filters ...FrameFilter) *Stacktrace {
	pcs := make([]uintptr, 100)
	n := runtime.Callers(1, pcs)

	if n == 0 {
		return nil
	}
	frames := extractFrames(pcs[:n])

	// default filter
	frames = filterFrames(frames, func(f Frame) bool {
		return f.Package == "runtime" || f.Package == "testing" ||
			strings.HasSuffix(f.Package, "/herror")
	})

	for _, filter := range filters {
		frames = filterFrames(frames, filter)
	}

	stacktrace := Stacktrace{
		Frames: frames,
	}

	return &stacktrace
}

// NewFrame assembles a stacktrace frame out of `runtime.Frame`.
func NewFrame(f runtime.Frame) Frame {
	abspath := unknown
	filename := unknown
	if f.File != "" {
		abspath = f.File
		_, filename = filepath.Split(f.File)
	}

	function := unknown
	pkgname := unknown
	typer := ""
	if f.Function != "" {
		pkgname, typer, function = deconstructFunctionName(f.Function)
	}

	inApp := func() bool {
		out := strings.HasPrefix(abspath, build.Default.GOROOT) ||
			strings.Contains(pkgname, "vendor")
		return !out
	}()

	return Frame{
		AbsPath:  abspath,
		Filename: filename,
		Line:     f.Line,
		Package:  pkgname,
		Type:     typer,
		Function: function,
		InApp:    inApp,
	}
}

func filterFrames(frames []Frame, filter FrameFilter) []Frame {
	filtered := make([]Frame, 0, len(frames))

	for _, frame := range frames {
		if filter(frame) {
			continue
		}
		filtered = append(filtered, frame)
	}

	return filtered
}

func extractFrames(pcs []uintptr) []Frame {
	frames := make([]Frame, 0, len(pcs))
	callersFrames := runtime.CallersFrames(pcs)

	for {
		callerFrame, more := callersFrames.Next()
		frames = append([]Frame{
			NewFrame(callerFrame),
		}, frames...)

		if !more {
			break
		}
	}

	return frames
}

func deconstructFunctionName(name string) (pkg string, typer string, function string) {
	if i := strings.LastIndex(name, "/"); i != -1 {
		pkg = name[:i]
		function = name[i+1:]

		if d := strings.Index(function, "."); d != -1 {
			pkg = fmt.Sprint(pkg, "/", function[:d])
			function = function[d+1:]
		}

		if o, c := strings.LastIndex(name, ".("), strings.LastIndex(name, ")."); o != -1 && c != -1 {
			pkg = name[:o]
			function = name[c+2:]

			typer = name[o+2 : c]
			if i := strings.Index(typer, "*"); i != -1 {
				typer = typer[1:]
			}
		}
		return
	}

	if i := strings.LastIndex(name, "."); i != -1 {
		pkg = name[:i]
		function = name[i+1:]
	}

	return
}
