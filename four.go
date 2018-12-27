package main

import (
	"bytes"
	"encoding/binary"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
	"github.com/robfig/cron"
	"log"
)

func main() {
	portaudio.Initialize()
	var filename = "test2.mp3"
	cron := cron.New()
	cron.AddFunc("@every 1s", func() {
		if err := PlayThree(filename); err != nil {
			log.Fatal(err)
		}
	})
	//cron.AddFunc("@every 1s", func() {
	//	if err := PlayTwo(filename); err != nil {
	//		log.Fatal(err)
	//	}
	//})
	cron.Run()
}

func PlayThree(filename string) error {
	log.Print("sadfasdf\n")
	// new 一个解析器
	decoder, err := mpg123.NewDecoder("")
	if err != nil {
		return err
	}
	defer decoder.Delete()
	//打开文件
	if err := decoder.Open(filename); err != nil {
		return err
	}
	defer decoder.Close()
	//设定音乐格式
	rate, channels, _ := decoder.GetFormat()
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)
	// 缓冲
	out := make([]int16, 8196)
	// 打开 portaudio 流
	stream, err := portaudio.OpenDefaultStream(
		0,
		channels,
		float64(rate),
		len(out),
		&out,
	)
	if err != nil {
		return err
	}
	defer stream.Close()
	if err := stream.Start(); err != nil {
		return err
	}
	defer stream.Stop()

	for {
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err != nil {
			if err == mpg123.EOF {
				break
			}
			return err
		}
		if err = binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out); err != nil {

			return err
		}
		if err = stream.Write(); err != nil {
			return err
		}
	}
	return nil
}
