package api

import (
	"fmt"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

func SendNotification(title string, subtitle string, message string) error {
	note := gosxnotifier.NewNotification(message)
	note.Title = title
	note.Subtitle = subtitle
	note.Sound = gosxnotifier.Basso
	note.Group = BundleID
	note.Link = "https://tracker-suivi.apps.cic.gc.ca/en/login" // or BundleID like: com.apple.Terminal
	note.AppIcon = "assets/ic_canada.png"
	note.ContentImage = "assets/ic_canada.png"

	// Then, push the notification
	err := note.Push()

	fmt.Println(err)
	return err
}
