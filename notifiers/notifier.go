package notifiers

type Notifier interface {
	Open() bool
	Message(message string) bool
	Close() bool
}
