package notifier

import (
	notifylib "github.com/gen2brain/beeep"
)

// Notify ...
func Notify(title string, msg string) error {
	return notifylib.Notify(title, msg, "")
}
