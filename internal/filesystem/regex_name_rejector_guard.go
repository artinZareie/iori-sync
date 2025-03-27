package filesystem

import (
	"regexp"
)

// This struct is redundant due to existence of `^` in regex, however I've decided to implement it for the sake of
// consistency and readability. This way, we can have a clear distinction between acceptor and rejector guards.
type RegexNameRejectorGuard struct {
	Regex string
}

func (r *RegexNameRejectorGuard) Accept(file File) bool {
	return true
}

func (r *RegexNameRejectorGuard) Reject(file File) bool {
	re, err := regexp.Compile(r.Regex)
	if err != nil {
		return true
	}

	return re.MatchString(file.Info.Name())
}
