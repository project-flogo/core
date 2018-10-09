package managed

type Status string

const (
	StatusStarted Status = "Started"
	StatusStopped        = "Stopped"
	StatusFailed         = "Failed"
)

type StatusInfo struct {
	Name   string
	Status Status
	Error  error
}
