package validation

import "fmt"

var allowed = map[string]bool{"queued": true, "running": true, "hold": true, "completed": true}

func Status(value string) error {
	if !allowed[value] {
		return fmt.Errorf("status must be queued, running, hold, or completed")
	}
	return nil
}
