package tool

import (
	"fmt"
	"io"
	"time"
)

// formatBytes translates the C function `max5data`. It formats a byte count
// into a human-readable string with suffixes (k, M, G, etc.), fitting
// within a 5-character width.
func formatBytes(bytes int64) string {
	const (
		oneKB = 1024
		oneMB = 1024 * oneKB
		oneGB = 1024 * oneMB
		oneTB = 1024 * oneGB
		onePB = 1024 * oneTB
	)

	switch {
	case bytes < 100000:
		return fmt.Sprintf("%5d", bytes)
	case bytes < 10000*oneKB:
		return fmt.Sprintf("%4dk", bytes/oneKB)
	case bytes < 100*oneMB:
		return fmt.Sprintf("%2d.%dM", bytes/oneMB, (bytes%oneMB)/(oneMB/10))
	case bytes < 10000*oneMB:
		return fmt.Sprintf("%4dM", bytes/oneMB)
	case bytes < 100*oneGB:
		return fmt.Sprintf("%2d.%dG", bytes/oneGB, (bytes%oneGB)/(oneGB/10))
	case bytes < 10000*oneGB:
		return fmt.Sprintf("%4dG", bytes/oneGB)
	case bytes < 10000*oneTB:
		return fmt.Sprintf("%4dT", bytes/oneTB)
	default:
		return fmt.Sprintf("%4dP", bytes/onePB)
	}
}

// formatSeconds translates the C function `time2str`. It formats a duration
// in seconds into a human-readable string (HH:MM:SS or Xd Yh).
func formatSeconds(seconds int64) string {
	if seconds <= 0 {
		return "--:--:--"
	}

	h := seconds / 3600
	if h <= 99 {
		m := (seconds % 3600) / 60
		s := seconds % 60
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}

	d := seconds / 86400
	h = (seconds % 86400) / 3600
	if d <= 999 {
		return fmt.Sprintf("%3dd %02dh", d, h)
	}
	return fmt.Sprintf("%7dd", d)
}

// TransferStats holds the progress data for a single ongoing transfer.
type TransferStats struct {
	DLTotal, DLNow, ULTotal, ULNow int64
}

// ProgressBar renders a command-line progress meter.
// It encapsulates the state previously held in global static variables
// in tool_progress.c.
type ProgressBar struct {
	Writer     io.Writer
	StartTime  time.Time
	lastRender time.Time
}

// NewProgressBar creates a new progress bar instance.
func NewProgressBar(writer io.Writer) *ProgressBar {
	return &ProgressBar{
		Writer:    writer,
		StartTime: time.Now(),
	}
}

// Render displays the progress meter line. This translates the logic from
// the C function `progress_meter`. It should be called periodically.
func (p *ProgressBar) Render(stats []TransferStats, final bool) {
	now := time.Now()
	// Render at most twice per second, or if it's the final call.
	if !final && now.Sub(p.lastRender) < 500*time.Millisecond {
		return
	}
	p.lastRender = now

	var totalDL, totalUL, currentDL, currentUL int64
	for _, s := range stats {
		totalDL += s.DLTotal
		totalUL += s.ULTotal
		currentDL += s.DLNow
		currentUL += s.ULNow
	}

	dlPercent := "--"
	if totalDL > 0 {
		dlPercent = fmt.Sprintf("%2d", (currentDL*100)/totalDL)
	}

	ulPercent := "--"
	if totalUL > 0 {
		ulPercent = fmt.Sprintf("%2d", (currentUL*100)/totalUL)
	}

	timeSpent := int64(now.Sub(p.StartTime).Seconds())

	// Simplified speed calculation. The C version uses a sliding window.
	var speed int64
	if timeSpent > 0 {
		speed = (currentDL + currentUL) / timeSpent
	}

	timeLeft := int64(0)
	if speed > 0 && totalDL > 0 {
		timeLeft = (totalDL - currentDL) / speed
	}

	// Render the progress bar line. Using \r to return to the start of the line.
	fmt.Fprintf(p.Writer,
		"\rDL%% UL%%  Dled  Uled  Xfers  Live Total     Current  Left    Speed\n"+
			"%-3s %-3s %s %s %5d %5d  %s  %s  %s %s",
		dlPercent,
		ulPercent,
		formatBytes(currentDL),
		formatBytes(currentUL),
		len(stats), // Xfers
		len(stats), // Live (simplified)
		formatSeconds(0), // Total time (simplified)
		formatSeconds(timeSpent),
		formatSeconds(timeLeft),
		formatBytes(speed),
	)

	if final {
		fmt.Fprintln(p.Writer)
	}
}