package utils

import(
    "os"
    "strings"
)

const progressBarLength int    = 68

// PrintProgressBar prints a bar of the current status of the progress
// 'progress' is current status and 'complete' is the completed status
func PrintProgressBar(progress int64, complete int64) {
    amount := int((int64(progressBarLength-1) * progress) / complete)
    rest := (progressBarLength - 1) - amount
    bar := strings.Repeat("=", amount) + ">" + strings.Repeat(".", rest)
    _, err := os.Stdout.Write([]byte("Progress: [" + bar + "]\r"))
    Checkerr(err)
}

// PrintCompletedBar prints a completed bar and a message
func PrintCompletedBar() {
    _, err := os.Stdout.Write([]byte("Progress: [" + strings.Repeat("=", progressBarLength) + "]\r"))
    Checkerr(err)
    _, err = os.Stdout.Write([]byte("Completed \n"))
    Checkerr(err)
}
