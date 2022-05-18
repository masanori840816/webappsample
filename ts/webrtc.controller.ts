export class WebRtcController {
    private webcamStream: MediaStream|null = null; 
    private peerConnection: RTCPeerConnection|null = null;
    private answerSentEvent: ((data: RTCSessionDescriptionInit) => void)|null = null;
    private candidateSentEvent: ((data: RTCIceCandidate) => void)|null = null;
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
    public connect() {
        if(this.webcamStream == null) {
            console.error("Local video was null");
            return;
        }
        this.peerConnection = new RTCPeerConnection();
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
        for(const track of this.webcamStream.getTracks()) {
            this.peerConnection.addTrack(track);
        }
        this.peerConnection.onicecandidate = ev => {
            if (ev.candidate == null ||
                this.candidateSentEvent == null) {
              return;
            }
            this.candidateSentEvent(ev.candidate);
        };
    }
    public addEvents(answerSentEvent: (data: RTCSessionDescriptionInit) => void,
        candidateSentEvent: (data: RTCIceCandidate) => void) {
        this.answerSentEvent = answerSentEvent;
        this.candidateSentEvent = candidateSentEvent;
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
}