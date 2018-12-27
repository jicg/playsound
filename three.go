package main

import (
	"bytes"
	"encoding/binary"
	"github.com/gordonklaus/portaudio"
	"github.com/robfig/cron"
	"github.com/tosone/minimp3"
	"io"
	"io/ioutil"
	"log"
)

func main() {
	portaudio.Initialize()
	var filename = "test2.mp3"
	cron := cron.New()
	cron.AddFunc("@every 1m1s", func() {
		if err := PlayTwo(filename); err != nil {
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

func PlayTwo(filename string) error {
	var (
		err  error
		data []byte
		dec  *minimp3.Decoder
	)

	//读取文件
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	//用minimp3 解析 音乐
	if dec, data, err = minimp3.DecodeFull(bs); err != nil {
		return err
	}
	//当退出函数时，释放资源
	defer dec.Close()

	// 定义缓冲，portaudio Write 从缓冲里面获取数据。
	out := make([]int16, 8196)
	//用 portaudio 打开音频流
	stream, err := portaudio.OpenDefaultStream(
		0,
		dec.Channels,
		float64(dec.SampleRate),
		len(out),
		&out,
	)
	if err != nil {
		return err
	}
	//当退出函数，释放资源
	defer stream.Close()
	//开始视频处理
	if err := stream.Start(); err != nil {
		return err
	}
	// 退出函数时，释放资源
	defer stream.Stop()
	byteffer := bytes.NewBuffer(data)

	//循环 来读取minimp3解析出来的数据，并给portalaudio处理
	for {
		audio := make([]byte, 2*len(out))
		_, err := byteffer.Read(audio)
		// 有错误，就退出。io.EOF 说明文件已经结束
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		err = binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out)
		// 有错误，就退出。io.EOF 说明文件已经结束
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if err := stream.Write(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}
