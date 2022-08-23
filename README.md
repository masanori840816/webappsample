# webappsample

## Description
This is a sample to try WebRTC SFU.

I use [pion/webrtc](https://github.com/pion/webrtc) and refer [example-webrtc-applications/sfu-ws](https://github.com/pion/example-webrtc-applications/tree/master/sfu-ws).

For signaling and sending text messages, I use Sever-Sent Events.

## Environments
* Go ver.go1.19 windows/amd64
* Node.js ver.18.7.0

## My blog posts
* [[Go] Try Server-Sent Events](https://dev.to/masanori_msl/go-try-server-sent-events-19fh)
* [[Go] Try Pion/WebRTC with SSE](https://dev.to/masanori_msl/go-try-pionwebrtc-with-sse-582)
* [[Go][Pion/WebRTC] Closing chan and adding DataChannel](https://dev.to/masanori_msl/gopionwebrtc-closing-chan-and-adding-datachannel-f83)
* [[Pion/WebRTC] Enabling and disabling the video track](https://dev.to/masanori_msl/pionwebrtc-enabling-and-disabling-the-video-track-ma8)
* [[WebRTC][Web Audio API] Identify who is vocalizing](https://dev.to/masanori_msl/webrtcweb-audio-api-identify-who-is-vocalizing-3d0b)

## Environment variables
Before you run this project, you need add "WEBRTCAPP_URL={Your server address}" as a enviroment variable.

## License
MIT