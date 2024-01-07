package tui

import "fmt"

func validateString(data string) error {
	if len(data) < 1 {
		return fmt.Errorf("DATA %s MUST HAVE AT LEAST 1 CHARATER", data)
	}
	return nil
}
