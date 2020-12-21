package notifier

import (
	"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	notifylib "github.com/gen2brain/beeep"
)

// Notify ...
func Notify(title string, msg string) error {
	return notifylib.Notify(title, msg, "")
}

// NotifyWithSound ...
func NotifyWithSound(title string, msg string) error {
	err := Notify(title, msg)
	if err != nil {
		return err
	}

	err = PlaySound()
	if err != nil {
		return err
	}

	return nil
}

// PlaySound ...
func PlaySound() error {
	f, err := mp3File()
	if err != nil {
		return err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

func mp3File() (io.ReadCloser, error) {
	r := strings.TrimRight(strings.TrimLeft(Mp3(), "\r\n"), "\r\n")
	mp3Bytes, err := hex.DecodeString(r)
	if err != nil {
		return nil, err
	}

	f := ioutil.NopCloser(bytes.NewReader(mp3Bytes))
	return f, nil
}
