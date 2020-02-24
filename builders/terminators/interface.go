package terminators

// SearchTerminator ...
type SearchTerminator interface {
	IsSearchTerminated(pass int) bool
}
