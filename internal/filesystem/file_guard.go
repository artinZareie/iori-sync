package filesystem

// FileGuard is an interface for a valid file filter. Any struct implementing
// this interface can be used to filter files. The condition to accept a file
// is `Accept() && !Reject()`, where reject is defined for the sake of simplicity.
// To implement an acceptor, implement `Reject(f File) {return false}`.
// To implement a rejector, implement `Accept(f File) {return true}`.
type FileGuard interface {
	Accept(File) bool
	Reject(File) bool
}
