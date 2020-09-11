package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
	"github.com/povilasv/prommod"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gin-gonic/autotls"
	"log"

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
	//Init()
	r := gin.Default()
	if err := prometheus.Register(prommod.NewCollector("sfu_ws")); err != nil {
		panic(err)
	}

	port := flag.String("p", "8443", "https port")
	flag.Parse()

	//http.Handle("/metrics", promhttp.Handler())

	// Websocket handle func
	r.GET("/ws", func(c *gin.Context) {
		room(c.Writer, c.Request)
	})
	//http.HandleFunc("/ws", room)

	// Html handle func
	r.GET("/", func(context *gin.Context) {
		web(context.Writer,context.Request)
	})
	//http.HandleFunc("/", web)

	// Support https, so we can test by lan
	fmt.Println("Web listening :" + *port)
	//panic(http.ListenAndServeTLS(":"+*port, "cert.pem", "key.pem", nil))
	//panic(http.ListenAndServe("0.0.0.0:8080", nil))
	log.Fatal(autotls.Run(r, "includeamin.kelidiha.com"))
}
