package main

import (
	"fmt"
	"github.com/aynakeya/go-mpv"
	"github.com/inancgumus/screen"
	"log"
	"math/rand"
	"time"
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

func play_ad(category string) {
	// If the category present is not in the following array pick a random catagory
	categories := [13]string{
		"pythondev", "javadev", "jsdev", "cppdev", "rustdev", "csharpdev",
		"phpdev", "golangdev", "swiftdev", "kotlindev", "sysdev",
		"webdev", "securitydev", // Note: Trailing comma is optional but idiomatic
	} // Array size is 13

	// Check if the provided category exists in the array
	found := false
	for _, validCategory := range categories {
		if category == validCategory {
			found = true
			break // Found it, no need to check further
		}
	}
	// If the category was not found in the list
	if !found {
		//fmt.Printf("Category '%s' not recognized. Selecting a random category.\n", category)

		// Seed the random number generator.
		// In a real application, this should ideally be done once, e.g., in main() or init().
		// Using a local source to avoid interfering with the global rand state.
		source := rand.NewSource(time.Now().UnixNano())
		localRand := rand.New(source)

		// Pick a random index from 0 to len(categories)-1
		randomIndex := localRand.Intn(len(categories))

		// Assign the randomly selected category
		category = categories[randomIndex]
	}
	// Send to r.sock recive on s.sock
	localfile := "http://localhost:8080/" + category
	m := mpv.Create()
	c := eventListener(m)
	//log.Println("audio-client-name", m.SetOptionString("audio-client-name", "AynaMpvCore"))
	//log.Println("volume", m.SetOption("volume", mpv.FORMAT_INT64, 100))
	log.Println("terminal", m.SetOptionString("terminal", "no"))
	log.Println("vo", m.SetOptionString("vo", "tct"))
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
	screen.Clear()
	screen.MoveTopLeft()
}
