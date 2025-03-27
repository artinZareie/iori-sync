package filesystem

import (
	"regexp"
)

type RegexNameAcceptorGuard struct {
	Regex string
}

func (r *RegexNameAcceptorGuard) Accept(file File) bool {
	re, err := regexp.Compile(r.Regex)
	if err != nil {
		return false
	}

	return re.MatchString(file.Info.Name())
}

func (r *RegexNameAcceptorGuard) Reject(file File) bool {
	return false
}
