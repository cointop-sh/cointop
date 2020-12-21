package beep

// Ctrl allows for pausing a Streamer.
//
// Wrap a Streamer in a Ctrl.
//
//   ctrl := &beep.Ctrl{Streamer: s}
//
// Then, we can pause the streaming (this will cause Ctrl to stream silence).
//
//   ctrl.Paused = true
//
// To completely stop a Ctrl before the wrapped Streamer is drained, just set the wrapped Streamer
// to nil.
//
//   ctrl.Streamer = nil
//
// If you're playing a Streamer wrapped in a Ctrl through the speaker, you need to lock and unlock
// the speaker when modifying the Ctrl to avoid race conditions.
//
//   speaker.Play(ctrl)
//   // ...
//   speaker.Lock()
//   ctrl.Paused = true
//   speaker.Unlock()
type Ctrl struct {
	Streamer Streamer
	Paused   bool
}

// Stream streams the wrapped Streamer, if not nil. If the Streamer is nil, Ctrl acts as drained.
// When paused, Ctrl streams silence.
func (c *Ctrl) Stream(samples [][2]float64) (n int, ok bool) {
	if c.Streamer == nil {
		return 0, false
	}
	if c.Paused {
		for i := range samples {
			samples[i] = [2]float64{}
		}
		return len(samples), true
	}
	return c.Streamer.Stream(samples)
}

// Err returns the error of the wrapped Streamer, if not nil.
func (c *Ctrl) Err() error {
	if c.Streamer == nil {
		return nil
	}
	return c.Streamer.Err()
}
