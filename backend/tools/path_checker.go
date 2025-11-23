package tools

// PathChecker abstracts permission checks for filesystem paths.
// Implemented by memory.AllowedDirsStore.
type PathChecker interface {
    IsAllowed(path string) bool
}
