package tools

// All returns every tool the agent is allowed to call. Adding a tool means
// writing it in this package and appending it here — callers never assemble
// the list themselves.
func All() []Definition {
	return []Definition{
		ReadFile,
		ListFiles,
		EditFile,
		Bash,
	}
}
