package update

import (
	"context"
	"fmt"
	"time"

	"github.com/jbowes/whatsnew"
	"github.com/vikingo-project/vsat/shared"
)

func CheckNewVersion() {
	for {
		ctx := context.Background()
		fut := whatsnew.Check(ctx, &whatsnew.Options{
			Slug:    "vikingo-project/vsat",
			Version: shared.Version,
		})
		if v, _ := fut.Get(); v != "" {
			fmt.Printf("new version available: %s\n", v)
		}
		time.Sleep(30 * time.Minute)
	}
}
