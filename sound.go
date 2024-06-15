package main

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var moveSound *beep.Buffer

func loadSounds() {
	f, err := os.Open("./assets/clear.wav")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		return
	}
	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode file: %v", err)
		return
	}
	moveSound = beep.NewBuffer(format)
	moveSound.Append(streamer)
	streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func playClearSound() {
	done := make(chan bool)
	speaker.Play(beep.Seq(moveSound.Streamer(0, moveSound.Len()), beep.Callback(func() {
		done <- true
	})))
	<-done
}
