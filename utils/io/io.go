package utilsio

import "fmt"

const (
	KB float64 = 1024
	MB float64 = 1048576
	GB float64 = 1073741824
)

// Returns a human readable size from a size in bytes
//
/* GetHumanReadableSize(1024) // "1KO" */
func GetHumanReadableSize(size int) string {
	floatSize := float64(size)

	switch {
	case floatSize >= GB:
		return fmt.Sprintf("%.1fGB", floatSize/GB)
	case floatSize >= MB:
		return fmt.Sprintf("%.1fMB", floatSize/MB)
	case floatSize >= KB:
		return fmt.Sprintf("%.1fKB", floatSize/KB)
	default:
		return fmt.Sprintf("%.1fB", floatSize)
	}
}
