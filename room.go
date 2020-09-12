package main

import (
	"fmt"
	sio "github.com/googollee/go-socket.io"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"io"
	"sync/atomic"
	"time"
)

// Peer config
var peerConnectionConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
	SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
}

var (
	// Media engine
	media webrtc.MediaEngine
	//media MediaEngine

	// API object
	api *webrtc.API

	// Publisher Peer
	pubCount    int32
	pubReceiver *webrtc.PeerConnection

	server *sio.Server

	// Broadcast channels
	broadcastHub   = newHub()
	requestChan    chan string
	ConnChan       chan sio.Conn
	localTrackChan chan *webrtc.Track
)

const (
	rtcpPLIInterval = time.Second * 3
)

func ManageSocket() {

}

func room() {
	//media = webrtc.MediaEngine{}
	//media.RegisterDefaultCodecs()
	//media.RegisterDefaultCodecs()
	//api = webrtc.NewAPI(webrtc.WithMediaEngine(media))

	for {
		println("inja")
		msg := <-requestChan
		conn := <-ConnChan
		println("unja")

		if atomic.LoadInt32(&pubCount) == 0 {
			atomic.AddInt32(&pubCount, 1)
			offer := webrtc.SessionDescription{}

			Decode(msg, &offer)
			media = webrtc.MediaEngine{}
			err := media.PopulateFromSDP(offer)
			if err != nil {
				panic(err)
			}
			println("72")
			api = webrtc.NewAPI(webrtc.WithMediaEngine(media))
			pubReceiver, _ = api.NewPeerConnection(peerConnectionConfig)
			pubReceiver.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
				println(state.String())
			})
			// Listen for ICE Candidates from the remote peer
			//pubReceiver.AddICECandidate(remoteCandidate)
			pubReceiver.OnICECandidate(func(i *webrtc.ICECandidate) {

				if i != nil {
					go conn.Emit("ice", i.ToJSON())
					println(i.String())
				}
			})

			if _, err = pubReceiver.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
				panic(err)
			} else if _, err = pubReceiver.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
				panic(err)
			}
			println("79")
			localTrackChan = make(chan *webrtc.Track)

			pubReceiver.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
				localTrack, newTrackErr := pubReceiver.NewTrack(track.PayloadType(), track.SSRC(), "video", "pion")
				if newTrackErr != nil {
					panic(newTrackErr)
				}
				go func() {
					ticker := time.NewTicker(rtcpPLIInterval)
					for range ticker.C {
						if rtcpSendErr := pubReceiver.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: track.SSRC()}}); rtcpSendErr != nil {
							fmt.Println(rtcpSendErr)
						}
					}
				}()
				println(track.Kind())
				localTrackChan <- localTrack
				rtpBuf := make([]byte, 1400)
				for {
					i, readErr := track.Read(rtpBuf)
					if readErr != nil {
						panic(readErr)
					}

					// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
					if _, err := localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
						panic(err)
					}
				}
			})
			err = pubReceiver.SetRemoteDescription(offer)
			if err != nil {
				panic(err)
			}

			println(149)
			answer, err := pubReceiver.CreateAnswer(nil)
			checkError(err)
			println(153)

			checkError(pubReceiver.SetLocalDescription(answer))

			println(160)
			go conn.Emit("sdp", Encode(*pubReceiver.LocalDescription()))
			println("164")

			pubReceiver.OnDataChannel(func(d *webrtc.DataChannel) {
				d.OnMessage(func(msg webrtc.DataChannelMessage) {
					println("aminajaml")
					broadcastHub.broadcastChannel <- msg.Data
				})
			})
			println("172")

		} else {
			println("155")
			// Create a new PeerConnection
			subSender, err := api.NewPeerConnection(peerConnectionConfig)
			//subSender, err := NewPeerConnection(peerConnectionConfig)
			checkError(err)
			println("159")
			subSender.OnICECandidate(func(i *webrtc.ICECandidate) {

				if i != nil {
					go conn.Emit("ice", i.ToJSON())
					println(i.String())
				}
			})

			// Register data channel creation handling
			subSender.OnDataChannel(func(d *webrtc.DataChannel) {
				broadcastHub.addListener(d)
			})
			println("166")

			_, err = subSender.AddTrack(<-localTrackChan)
			if err != nil {
				panic(err)
			}
			println("196")
			recvOnlyOffer := webrtc.SessionDescription{}
			Decode(msg, &recvOnlyOffer)
			checkError(subSender.SetRemoteDescription(recvOnlyOffer))
			println("205")

			answer, err := subSender.CreateAnswer(nil)
			checkError(err)
			println("211")

			checkError(subSender.SetLocalDescription(answer))
			println("215")

			// Send server sdp to subscriber
			println("hre")

			go conn.Emit("sdp", Encode(*subSender.LocalDescription()))

		}
	}
}
