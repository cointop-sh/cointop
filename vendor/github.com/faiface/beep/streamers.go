package beep

// Silence returns a Streamer which streams num samples of silence. If num is negative, silence is
// streamed forever.
func Silence(num int) Streamer {
	return StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		if num == 0 {
			return 0, false
		}
		if 0 < num && num < len(samples) {
			samples = samples[:num]
		}
		for i := range samples {
			samples[i] = [2]float64{}
		}
		if num > 0 {
			num -= len(samples)
		}
		return len(samples), true
	})
}

// Callback returns a Streamer, which does not stream any samples, but instead calls f the first
// time its Stream method is called.
func Callback(f func()) Streamer {
	return StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		if f != nil {
			f()
			f = nil
		}
		return 0, false
	})
}

// Iterate returns a Streamer which successively streams Streamers obtains by calling the provided g
// function. The streaming stops when g returns nil.
//
// Iterate does not propagate errors from the generated Streamers.
func Iterate(g func() Streamer) Streamer {
	var (
		s     Streamer
		first = true
	)
	return StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		if first {
			s = g()
			first = false
		}
		if s == nil {
			return 0, false
		}
		for len(samples) > 0 {
			if s == nil {
				break
			}
			sn, sok := s.Stream(samples)
			if !sok {
				s = g()
			}
			samples = samples[sn:]
			n += sn
		}
		return n, true
	})
}
