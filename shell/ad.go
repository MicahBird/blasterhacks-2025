package main

import (
	"fmt"
	"github.com/aynakeya/go-mpv"
	"log"
)

type AdType string

const (
	PROGRAMMER = AdType("programmer")
	SOY_DEV    = AdType("SOY_DEV")
)

func eventListener(m *mpv.Mpv) chan *mpv.Event {
	c := make(chan *mpv.Event)
	go func() {
		for {
			e := m.WaitEvent(1)
			c <- e
		}
	}()
	return c
}

func play_ad(ad_type AdType) {
	localfile := "/home/user/" + ad_type + ".mp4"
	m := mpv.Create()
	c := eventListener(m)
	//log.Println("audio-client-name", m.SetOptionString("audio-client-name", "AynaMpvCore"))
	//log.Println("volume", m.SetOption("volume", mpv.FORMAT_INT64, 100))
	log.Println("terminal", m.SetOptionString("terminal", "no"))
	log.Println("vo", m.SetOptionString("vo", "caca"))
	//log.Println("set ao", m.SetPropertyString("audio-device", "pulse/alsa_output.pci-0000_75_00.6.analog-stereo"))

	//log.Println("video", m.SetOption("video", mpv.FORMAT_STRING, "no"))
	//log.Println("vo=null", m.SetOption("vo", mpv.FORMAT_STRING, "null"))
	//log.Println("vo=null", m.SetOptionString("vo", "null"))
	//log.Println("vo=null", m.SetPropertyString("vo", "null"))
	//log.Println("vid", m.SetOption("vid", mpv.FORMAT_STRING, "no"))

	err := m.Initialize()

	if err != nil {
		log.Println("Mpv init:", err.Error())
		return
	}
	//Set video file
	log.Println("loadfile", m.Command([]string{"loadfile", string(localfile)}))

	// getting log messages
	//m.RequestLogMessages(mpv.LOG_LEVEL_INFO)

	//m.ObserveProperty(1, "time-pos", mpv.FORMAT_NODE)
	m.ObserveProperty(1, "time-pos", mpv.FORMAT_STRING)

	fmt.Println(123)
	//fmt.Println(m.GetProperty("volume", mpv.FORMAT_NODE))
	//fmt.Println(m.GetProperty("ao-volume", mpv.FORMAT_NODE))

	for {
		e := <-c
		//log.Println(e)
		//if e.EventId == mpv.EVENT_LOG_MESSAGE {
		//	fmt.Println(e.LogMessage())
		//}
		//if e.EventId == mpv.EVENT_PROPERTY_CHANGE {
		//	fmt.Println(e.Property())
		//}
		if e.EventId == mpv.EVENT_END_FILE {
			fmt.Println("Thanks for watching!")
			break
		}

	}
	m.TerminateDestroy()
}
