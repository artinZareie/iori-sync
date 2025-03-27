package filesystem

import (
	"regexp"
)

// This struct is redundant due to existence of `^` in regex, however I've decided to implement it for the sake of
// consistency and readability. This way, we can have a clear distinction between acceptor and rejector guards.
type RegexPathRejectorGuard struct {
	Regex string
}

func (r *RegexPathRejectorGuard) Accept(file File) bool {
	return true
}

func (r *RegexPathRejectorGuard) Reject(file File) bool {
	re, err := regexp.Compile(r.Regex)
	if err != nil {
		return true
	}

	return re.MatchString(file.Path)
}
