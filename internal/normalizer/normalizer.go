package normalizer

type normalizer struct {
	filePath string
	tui      bool
	debug    bool
}

func WithFilePath(filePath string) func(*normalizer) {
	return func(a *normalizer) {
		a.filePath = filePath
	}
}

func DisableTui(disable bool) func(*normalizer) {
	return func(a *normalizer) {
		a.tui = !disable
	}
}

func WithDebug(debug bool) func(*normalizer) {
	return func(a *normalizer) {
		a.debug = debug
	}
}
