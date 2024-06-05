package autorepair

type NotificationLevel string

const (
	Info    NotificationLevel = "info"
	Success NotificationLevel = "success"
	Warning NotificationLevel = "warning"
	Fail    NotificationLevel = "fail"
)

type Notification struct {
	Level   NotificationLevel
	Message string
}
