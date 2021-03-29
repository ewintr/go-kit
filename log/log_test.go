package log_test

type testWriter struct {
	LogLines []string
}

func (tw *testWriter) Write(p []byte) (int, error) {
	tw.LogLines = append(tw.LogLines, string(p))

	return len(p), nil
}

func (tw *testWriter) Flush() {
	tw.LogLines = []string{}
}
