package notifiers

type Notifier interface {
	Open(configPath string) bool
	Message(message string) bool
	Close() bool
}
