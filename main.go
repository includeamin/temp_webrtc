package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3"
	"log"
	"net/http"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func NewRTPH264Codec(payloadType uint8, clockrate uint32, profileLevelId string) *webrtc.RTPCodec {
	return webrtc.NewRTPCodec(webrtc.RTPCodecTypeVideo,
		webrtc.H264,
		clockrate,
		0,
		fmt.Sprintf("level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=%s", profileLevelId),
		payloadType,
		&codecs.H264Payloader{})
}

func NewPeerConnection(configuration webrtc.Configuration) (*webrtc.PeerConnection, error) {

	media = webrtc.MediaEngine{}
	media.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))
	media.RegisterCodec(NewRTPH264Codec(webrtc.DefaultPayloadTypeH264, 90000, "42e01f"))
	api := webrtc.NewAPI(webrtc.WithMediaEngine(media))
	return api.NewPeerConnection(configuration)
}
func Init() {
	// Generate pem file for https
	//genPem()
	//
	//Create a MediaEngine object to configure the supported codec
	media = webrtc.MediaEngine{}
	media.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))
	media.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	//media = sfu.MediaEngine{}

	// Setup the codecs you want to use.
	//media.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	////media.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeH264, 90000))
	//media.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))
	//
	//
	////Create the API object with the MediaEngine
	//api = webrtc.NewAPI(webrtc.WithMediaEngine(media))
	api = webrtc.NewAPI(webrtc.WithMediaEngine(media))

}
func InitSocketIo() {




}

func main() {
	//Init()
	//InitSocketIo()
	t, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server = t
	requestChan = make(chan string,2)
	ConnChan = make(chan socketio.Conn,2)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "sdp", func(s socketio.Conn, msg string) {
		requestChan <- msg
		ConnChan <- s
		println("after")
	})
	go server.Serve()
	defer server.Close()
	go room()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	http.HandleFunc("/index",web)
	log.Println("Serving at localhost:8000...")
	//log.Fatal(http.ListenAndServe(":8000", nil))
	http.ListenAndServeTLS(":8443","cert.pem","key.pem",nil)





	////r := gin.Default()
	//if err := prometheus.Register(prommod.NewCollector("sfu_ws")); err != nil {
	//	panic(err)
	//}
	//
	//port := flag.String("p", "8443", "https port")
	//flag.Parse()
	//
	//http.Handle("/metrics", promhttp.Handler())
	//
	//// Websocket handle func
	////r.GET("/ws", func(c *gin.Context) {
	////	room(c.Writer, c.Request)
	////})
	//http.HandleFunc("/ws", room)
	//
	//// Html handle func
	/////	r.GET("/", func(context *gin.Context) {
	////		web(context.Writer, context.Request)
	////	})
	//http.HandleFunc("/", web)
	////r.GET("/ping", func(c *gin.Context) {
	////	c.String(200, "pong")
	////})
	//
	//// Support https, so we can test by lan
	//fmt.Println("Web listening :" + *port)
	//panic(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))
	////panic(http.ListenAndServe("0.0.0.0:8080", nil))
	////log.Fatal(autotls.Run(r, "includeamin.kelidiha.com"))
	////r.RunTLS(":8443", "./cert.pem", "./key.pem")
	////r.Run("0.0.0.0:8080")
}
