package support

type NeedsCleanup interface {
	Cleanup() error
}