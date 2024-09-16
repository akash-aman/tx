package constants

import "fmt"

const (
	CALCULATE_COUNT_MESSAGE = "Calculating template count ..."
)

func TEMPLATE_COUNT_MESSAGE(count int) string {
	return fmt.Sprintf("Got %d templates ...", count)
}
