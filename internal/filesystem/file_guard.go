package filesystem


type FileGuard interface {
	Accept(File) bool
	Reject(File) bool
}
