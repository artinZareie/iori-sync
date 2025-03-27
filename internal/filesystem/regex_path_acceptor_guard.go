package filesystem

import (
	"regexp"
)

type RegexPathAcceptorGuard struct {
	Regex string
}

func (r *RegexPathAcceptorGuard) Accept(file File) bool {
	re, err := regexp.Compile(r.Regex)
	if err != nil {
		return false
	}

	return re.MatchString(file.Path)
}

func (r *RegexPathAcceptorGuard) Reject(File) bool {
	return false
}
