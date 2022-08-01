import * as dataChannel from "./dataChannels"

export class WebRtcController {
    private webcamStream: MediaStream|null = null; 
    private peerConnection: RTCPeerConnection|null = null;
    private dataChannels: dataChannel.DataChannel[] = [];
    private answerSentEvent: ((data: RTCSessionDescriptionInit) => void)|null = null;
    private candidateSentEvent: ((data: RTCIceCandidate) => void)|null = null;
    private dataChannelMessageEvent: ((data: string|Uint8Array) => void)|null = null;
    private connectionUpdatedEvent: (() => void)|null = null;
    private localVideo: HTMLVideoElement;
    public constructor() {
        this.localVideo = document.getElementById("local_video") as HTMLVideoElement;
    }
    public init() {
        this.localVideo.addEventListener("canplay", () => {
            const width = 320;
            const height = this.localVideo.videoHeight / (this.localVideo.videoWidth/width);          
            this.localVideo.setAttribute("width", width.toString());
            this.localVideo.setAttribute("height", height.toString());
          }, false);
          navigator.mediaDevices.getUserMedia({ video: this.checkParams("video"), audio: true })
            .then(stream => {
                this.webcamStream = stream;
            });
    }
    public addEvents(answerSentEvent: (data: RTCSessionDescriptionInit) => void,
        candidateSentEvent: (data: RTCIceCandidate) => void,
        dataChannelMessageEvent: (data: string|Uint8Array) => void,
        connectionUpdatedEvent: () => void) {
        this.answerSentEvent = answerSentEvent;
        this.candidateSentEvent = candidateSentEvent;
        this.dataChannelMessageEvent = dataChannelMessageEvent;
        this.connectionUpdatedEvent = connectionUpdatedEvent;
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

        this.peerConnection.onconnectionstatechange = (ev) => {
            console.log(ev);
        };
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
            this.peerConnection.addTrack(track, this.webcamStream);
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
    public switchLocalVideoUsage(used: boolean): void {
        if(this.peerConnection == null ||
            this.webcamStream == null) {
            return;
        }
        const senders = this.peerConnection.getSenders();   
        console.log(senders);
        console.log("trans");
        console.log(this.peerConnection.getTransceivers());
        console.log("rec");
        console.log(this.peerConnection.getReceivers());
        
        const tracks = this.webcamStream.getVideoTracks();
        console.log("tracks");
        console.log(tracks);
        
        if(used) {
            if(tracks.length > 0) {
                console.log("2nd");
                
                /*let sender: RTCRtpSender|null = null;
                for(const s of senders) {
                    if(s.track?.kind === "video") {
                        sender = s;
                        console.log("Removeold");
                    }
                }*/
                if(tracks[0] == null) {
                    return;
                }
                for(const s of this.peerConnection.getSenders()) {
                    if(s.track == null || s.track.kind === "video") {
                        s.replaceTrack(tracks[0]);
                    }
                }
                /*for(const v of tracks) {
                    console.log("2nd add track");
                    console.log(v);
                    
                    this.peerConnection.addTrack(v, this.webcamStream);
                }*/
                console.log("result 2nd");
                
                console.log(this.peerConnection.getSenders());
                
        console.log(this.peerConnection.getTransceivers());
                if(this.connectionUpdatedEvent != null) {
                    this.connectionUpdatedEvent();
                }
                /*
                if(this.connectionUpdatedEvent != null) {
                    this.connectionUpdatedEvent();
                }*/
            } else {
                navigator.mediaDevices.getUserMedia({ video: true })
                .then(stream => {
                    const newVideoTracks = stream.getVideoTracks();
                    if(this.webcamStream == null ||
                        this.peerConnection == null ||
                        newVideoTracks.length <= 0) {
                        return;
                    }
                    this.localVideo.srcObject = stream;
                    this.localVideo.play();
                    for(const v of newVideoTracks) {
                        console.log("add new Videos");
                        
                        this.webcamStream.addTrack(v);
                        this.peerConnection.addTrack(v, this.webcamStream);
                    
                    }
                    
                    /*console.log("add track2");
                    this.webcamStream.addTrack(newVideoTracks[0]);
                    if(sender == null) {
                        this.peerConnection.addTrack(newVideoTracks[0], this.webcamStream);
                    } else {
                        console.log("add track34");
                                 
                        this.peerConnection.removeTrack(sender);
                        //const tranceiver = this.peerConnection.addTransceiver(newVideoTracks[0]);
                        //tranceiver.direction = "sendonly";
                        this.peerConnection.addTrack(newVideoTracks[0], this.webcamStream);
                    }
                    console.log(newVideoTracks[0]);*/

                //const tranceiver = this.peerConnection?.addTransceiver(newVideoTracks[0]);
                  ///  tranceiver!.direction = "sendonly";

                  console.log("result 1st");
                
                  console.log(this.peerConnection.getSenders());
                  
          console.log(this.peerConnection.getTransceivers());
                  if(this.connectionUpdatedEvent != null) {
                    this.connectionUpdatedEvent();
                }
                    
                    
                }); 
            }
        } else {
            if(senders.length > 0) {
                    
            console.log("Not Use2");
                //this.localVideo.pause();
                
                // tracks[0].stop();
                let sender: RTCRtpSender|null = null;
                for(const s of senders) {
                    if(s.track?.kind === "video") {
                        console.log("RemoveTrack");
                        
                        sender = s;
                        this.peerConnection.removeTrack(s);
                 //       s.track.stop();
                    }
                }
                if(sender != null) {
                    //this.peerConnection.removeTrack(sender);
                    //console.log(this.peerConnection);
                    
                    //const tranceiver = this.peerConnection.addTransceiver(tracks[0]);
                    
                }
                /*for(const t of tracks) {
                    console.log("RemoveTrack");
                    
                    this.webcamStream.removeTrack(t);
                
                }*/
                if(this.connectionUpdatedEvent != null) {
                    this.connectionUpdatedEvent();
                }
            }            
        }
    }
    private checkParams(key: string): boolean {
        const splittedUrls = location.href.split("/");
        const params = new URLSearchParams(splittedUrls[splittedUrls.length - 1]);
        const target = params.get(key);
        if(target == null || target.length <= 0) {
            return false;
        }
        return target === "true";
    }
}