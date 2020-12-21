package beep

// Streamer is able to stream a finite or infinite sequence of audio samples.
type Streamer interface {
	// Stream copies at most len(samples) next audio samples to the samples slice.
	//
	// The sample rate of the samples is unspecified in general, but should be specified for
	// each concrete Streamer.
	//
	// The value at samples[i][0] is the value of the left channel of the i-th sample.
	// Similarly, samples[i][1] is the value of the right channel of the i-th sample.
	//
	// Stream returns the number of streamed samples. If the Streamer is drained and no more
	// samples will be produced, it returns 0 and false. Stream must not touch any samples
	// outside samples[:n].
	//
	// There are 3 valid return pattterns of the Stream method:
	//
	//   1. n == len(samples) && ok
	//
	// Stream streamed all of the requested samples. Cases 1, 2 and 3 may occur in the following
	// calls.
	//
	//   2. 0 < n && n < len(samples) && ok
	//
	// Stream streamed n samples and drained the Streamer. Only case 3 may occur in the
	// following calls.
	//
	//   3. n == 0 && !ok
	//
	// The Streamer is drained and no more samples will come. If Err returns a non-nil error, only
	// this case is valid. Only this case may occur in the following calls.
	Stream(samples [][2]float64) (n int, ok bool)

	// Err returns an error which occurred during streaming. If no error occurred, nil is
	// returned.
	//
	// When an error occurs, Streamer must become drained and Stream must return 0, false
	// forever.
	//
	// The reason why Stream doesn't return an error is that it dramatically simplifies
	// programming with Streamer. It's not very important to catch the error right when it
	// happens.
	Err() error
}

// StreamSeeker is a finite duration Streamer which supports seeking to an arbitrary position.
type StreamSeeker interface {
	Streamer

	// Duration returns the total number of samples of the Streamer.
	Len() int

	// Position returns the current position of the Streamer. This value is between 0 and the
	// total length.
	Position() int

	// Seek sets the position of the Streamer to the provided value.
	//
	// If an error occurs during seeking, the position remains unchanged. This error will not be
	// returned through the Streamer's Err method.
	Seek(p int) error
}

// StreamCloser is a Streamer streaming from a resource which needs to be released, such as a file
// or a network connection.
type StreamCloser interface {
	Streamer

	// Close closes the Streamer and releases it's resources. Streamer will no longer stream any
	// samples.
	Close() error
}

// StreamSeekCloser is a union of StreamSeeker and StreamCloser.
type StreamSeekCloser interface {
	Streamer
	Len() int
	Position() int
	Seek(p int) error
	Close() error
}

// StreamerFunc is a Streamer created by simply wrapping a streaming function (usually a closure,
// which encloses a time tracking variable). This sometimes simplifies creating new streamers.
//
// Example:
//
//   noise := StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
//       for i := range samples {
//           samples[i][0] = rand.Float64()*2 - 1
//           samples[i][1] = rand.Float64()*2 - 1
//       }
//       return len(samples), true
//   })
type StreamerFunc func(samples [][2]float64) (n int, ok bool)

// Stream calls the wrapped streaming function.
func (sf StreamerFunc) Stream(samples [][2]float64) (n int, ok bool) {
	return sf(samples)
}

// Err always returns nil.
func (sf StreamerFunc) Err() error {
	return nil
}
