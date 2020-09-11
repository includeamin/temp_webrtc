package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/pion/webrtc/v3"
	"github.com/povilasv/prommod"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Init() {
	// Generate pem file for https
	genPem()

	// Create a MediaEngine object to configure the supported codec
	media = webrtc.MediaEngine{}
	//media = sfu.MediaEngine{}

	// Setup the codecs you want to use.
	media.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	media.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	// Create the API object with the MediaEngine
	api = webrtc.NewAPI(webrtc.WithMediaEngine(media))

}

func main() {
	Init()
	if err := prometheus.Register(prommod.NewCollector("sfu_ws")); err != nil {
		panic(err)
	}

	port := flag.String("p", "8443", "https port")
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())

	// Websocket handle func
	http.HandleFunc("/ws", room)

	// Html handle func
	http.HandleFunc("/", web)

	// Support https, so we can test by lan
	fmt.Println("Web listening :" + *port)
	//panic(http.ListenAndServeTLS(":"+*port, "cert.pem", "key.pem", nil))
	panic(http.ListenAndServe("0.0.0.0:8080", nil))
}
