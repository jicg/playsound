package main

import (
	"io/ioutil"

	"github.com/hajimehoshi/oto"
	"github.com/robfig/cron"
	"github.com/tosone/minimp3"
	"log"
)

type MyJOb struct {
}

func (this *MyJOb) Run() {
}

func main() {
	cron := cron.New()
	cron.AddFunc("", func() {
	})
	cron.AddJob("", &MyJOb{})
	cron.AddFunc("@every 1s", func() {
		if err := PlayOne("test.mp3"); err != nil {
			log.Fatal(err)
		}
	})
	cron.AddFunc("@every 1s", func() {
		if err := PlayOne("test.mp3"); err != nil {
			log.Fatal(err)
		}
	})
	cron.AddFunc("@every 1s", func() {
		if err := PlayOne("test.mp3"); err != nil {
			log.Fatal(err)
		}
	})
	cron.AddFunc("@every 1s", func() {
		if err := PlayOne("test.mp3"); err != nil {
			log.Fatal(err)
		}
	})
	cron.AddFunc("@every 1s", func() {
		if err := PlayOne("test.mp3"); err != nil {
			log.Fatal(err)
		}
	})
	cron.Run()
}
func PlayOne(name string) error {
	var file, err = ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	dec, data, err := minimp3.DecodeFull(file)
	if err != nil {
		return err
	}
	defer dec.Close()
	player, err := oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024)
	if err != nil {
		return err
	}
	defer player.Close()
	_, err = player.Write(data)
	if err != nil {
		return err
	}
	return nil
}
