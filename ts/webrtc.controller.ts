import * as dataChannel from "./dataChannels"

export class WebRtcController {
    private webcamStream: MediaStream|null = null; 
    private peerConnection: RTCPeerConnection|null = null;
    private dataChannels: dataChannel.DataChannel[] = [];
    private answerSentEvent: ((data: RTCSessionDescriptionInit) => void)|null = null;
    private candidateSentEvent: ((data: RTCIceCandidate) => void)|null = null;
    private dataChannelMessageEvent: ((data: string|Uint8Array) => void)|null = null;
    public init() {
        const localVideo = document.getElementById("local_video") as HTMLVideoElement;
        localVideo.addEventListener("canplay", () => {
            const width = 320;
            const height = localVideo.videoHeight / (localVideo.videoWidth/width);          
            localVideo.setAttribute("width", width.toString());
            localVideo.setAttribute("height", height.toString());
          }, false);
        navigator.mediaDevices.getUserMedia({ video: true, audio: true })
          .then(stream => {
              localVideo.srcObject = stream;
              localVideo.play();
              this.webcamStream = stream;
          })
          .catch(err => console.error(`An error occurred: ${err}`));
    }
    public addEvents(answerSentEvent: (data: RTCSessionDescriptionInit) => void,
        candidateSentEvent: (data: RTCIceCandidate) => void,
        dataChannelMessageEvent: (data: string|Uint8Array) => void) {
        this.answerSentEvent = answerSentEvent;
        this.candidateSentEvent = candidateSentEvent;
        this.dataChannelMessageEvent = dataChannelMessageEvent;
    }
    public handleOffer(data: RTCSessionDescription|null|undefined) {
        if(this.peerConnection == null ||
                data == null) {
            console.error("PeerConnection|SDP was null");
            return;
        }
        this.peerConnection.setRemoteDescription(data);
        this.peerConnection.createAnswer()
            .then(answer => {
                if(this.peerConnection != null) {
                    this.peerConnection.setLocalDescription(answer);
                }
                if(this.answerSentEvent != null) {
                    this.answerSentEvent(answer);
                }
            });
    }
    public handleCandidate(data: RTCIceCandidate|null|undefined) {
        if(this.peerConnection == null ||
            data == null) {
            console.error("PeerConnection|Candidate was null");
            return;
        }
        this.peerConnection.addIceCandidate(data);
    }
    public sendTextDataChannel(value: string) {
        if(this.peerConnection == null ||
            this.peerConnection.connectionState !== "connected") {
            return;
        }
        const target = this.dataChannels.find(c => c.dataChannel.id == 20);
        if(target == null) {
            return;
        }
        target.dataChannel.send(value);
    }
    public connect() {
        if(this.webcamStream == null) {
            console.error("Local video was null");
            return;
        }
        this.peerConnection = new RTCPeerConnection({
            iceServers: [{
                urls: `stun:stun.l.google.com:19302`,  // A STUN server              
            }]
        });

        this.peerConnection.onconnectionstatechange = (ev) => console.log(ev);
        
        this.peerConnection.ontrack = (ev) => {
            if (ev.track.kind === "audio" ||
                ev.streams[0] == null) {
              return;
            }    
            const remoteVideo = document.createElement("video");
            remoteVideo.srcObject = ev.streams[0];
            remoteVideo.autoplay = true;
            remoteVideo.controls = true;
            const videoArea = document.getElementById("remote_video_area") as HTMLElement;
            videoArea.appendChild(remoteVideo);
    
            ev.track.onmute = () => {
                remoteVideo.play();
            };
    
            ev.streams[0].onremovetrack = () => {
              if (remoteVideo.parentNode) {
                remoteVideo.parentNode.removeChild(remoteVideo);
              }
            };
          };
        this.webcamStream.getTracks().forEach(track => {
            if(this.peerConnection == null ||
                this.webcamStream == null) {
                return;
            }
            this.peerConnection.addTrack(track, this.webcamStream)
        });
        this.peerConnection.onicecandidate = ev => {
            if (ev.candidate == null ||
                this.candidateSentEvent == null) {
              return;
            }
            this.candidateSentEvent(ev.candidate);
        };        
        this.dataChannels.push(
            dataChannel.createTextDataChannel("sample3", 20, this.peerConnection,
            (message) => {
                if(this.dataChannelMessageEvent != null) {
                    this.dataChannelMessageEvent(message);
                }
            }));
        this.dataChannels.push(
            dataChannel.createTextDataChannel("sample2", 21, this.peerConnection,
            (message) => {
                if(this.dataChannelMessageEvent != null) {
                    this.dataChannelMessageEvent(message);
                }
            }));
    }
}