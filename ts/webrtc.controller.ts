import * as dataChannel from "./dataChannels"

export class WebRtcController {
    private webcamStream: MediaStream | null = null;
    private peerConnection: RTCPeerConnection | null = null;
    private dataChannels: dataChannel.DataChannel[] = [];
    private answerSentEvent: ((data: RTCSessionDescriptionInit) => void) | null = null;
    private candidateSentEvent: ((data: RTCIceCandidate) => void) | null = null;
    private dataChannelMessageEvent: ((data: string | Uint8Array) => void) | null = null;
    private connectionUpdatedEvent: (() => void) | null = null;
    private streamReceivedEvent: ((stream: MediaStream, kind: "video"|"audio") => void)|null = null;
    private streamRemovedEvent: ((id: string, kind: "video"|"audio") => void)|null = null;
    private localVideo: HTMLVideoElement;
    public constructor() {
        this.localVideo = document.getElementById("local_video") as HTMLVideoElement;
    }
    public init(videoUsed: boolean) {
        this.localVideo.addEventListener("canplay", () => {
            const width = 320;
            const height = this.localVideo.videoHeight / (this.localVideo.videoWidth / width);
            this.localVideo.setAttribute("width", width.toString());
            this.localVideo.setAttribute("height", height.toString());
        }, false);
        navigator.mediaDevices.getUserMedia({ video: videoUsed, audio: true })
            .then(stream => {
                this.webcamStream = stream;
            });
    }
    public addEvents(answerSentEvent: (data: RTCSessionDescriptionInit) => void,
        candidateSentEvent: (data: RTCIceCandidate) => void,
        dataChannelMessageEvent: (data: string | Uint8Array) => void,
        connectionUpdatedEvent: () => void,
        streamReceivedEvent: ((stream: MediaStream, kind: "video"|"audio") => void),
        streamRemovedEvent: ((id: string, kind: "video"|"audio") => void)) {
        this.answerSentEvent = answerSentEvent;
        this.candidateSentEvent = candidateSentEvent;
        this.dataChannelMessageEvent = dataChannelMessageEvent;
        this.connectionUpdatedEvent = connectionUpdatedEvent;
        this.streamReceivedEvent = streamReceivedEvent;
        this.streamRemovedEvent = streamRemovedEvent;
    }
    public handleOffer(data: RTCSessionDescription | null | undefined) {
        if (this.peerConnection == null ||
            data == null) {
            console.error("PeerConnection|SDP was null");
            return;
        }
        this.peerConnection.setRemoteDescription(data);
        this.peerConnection.createAnswer()
            .then(answer => {
                if (this.peerConnection != null) {
                    this.peerConnection.setLocalDescription(answer);
                }
                if (this.answerSentEvent != null) {
                    this.answerSentEvent(answer);
                }
            });
    }
    public handleCandidate(data: RTCIceCandidate | null | undefined) {
        if (this.peerConnection == null ||
            data == null) {
            console.error("PeerConnection|Candidate was null");
            return;
        }
        this.peerConnection.addIceCandidate(data);
    }
    public sendTextDataChannel(value: string) {
        if (this.peerConnection == null ||
            this.peerConnection.connectionState !== "connected") {
            return;
        }
        const target = this.dataChannels.find(c => c.dataChannel.id == 20);
        if (target == null) {
            return;
        }
        target.dataChannel.send(value);
    }
    public connect() {
        if (this.webcamStream == null) {
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
            if(this.streamReceivedEvent == null ||
                ev.streams[0] == null ||
                (ev.track.kind !== "audio" && ev.track.kind !== "video")) {
                return;
            }
            const stream = ev.streams[0];
            const kind = ev.track.kind;
            this.streamReceivedEvent(stream, ev.track.kind);

            stream.onremovetrack = () => {
                if (this.streamRemovedEvent != null) {
                    this.streamRemovedEvent(stream.id, kind);
                }
            };
        };
        this.webcamStream.getTracks().forEach(track => {
            if (this.peerConnection == null ||
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
                    if (this.dataChannelMessageEvent != null) {
                        this.dataChannelMessageEvent(message);
                    }
                }));
        this.dataChannels.push(
            dataChannel.createTextDataChannel("sample2", 21, this.peerConnection,
                (message) => {
                    if (this.dataChannelMessageEvent != null) {
                        this.dataChannelMessageEvent(message);
                    }
                }));
    }
    public switchLocalVideoUsage(used: boolean): void {
        if (this.peerConnection == null ||
            this.webcamStream == null) {
            return;
        }
        const tracks = this.webcamStream.getVideoTracks();
        if (used) {
            if (tracks.length > 0 &&
                tracks[0] != null) {
                this.replaceVideoTrack(this.peerConnection, tracks[0]);
            } else {
                this.addVideoTrack(this.peerConnection);
            }
        } else {
            this.removeVideoTrack(this.peerConnection);
        }
    }
    private addVideoTrack(peerConnection: RTCPeerConnection) {
        navigator.mediaDevices.getUserMedia({ video: true })
            .then(stream => {
                const newVideoTracks = stream.getVideoTracks();
                if (this.webcamStream == null ||
                    newVideoTracks.length <= 0) {
                    return;
                }
                this.localVideo.srcObject = stream;
                this.localVideo.play();
                for (const v of newVideoTracks) {
                    this.webcamStream.addTrack(v);
                    peerConnection.addTrack(v, this.webcamStream);
                }
                if (this.connectionUpdatedEvent != null) {
                    this.connectionUpdatedEvent();
                }
            });
    }
    private replaceVideoTrack(peerConnection: RTCPeerConnection, track: MediaStreamTrack) {
        this.localVideo.play();
        for (const s of peerConnection.getSenders()) {
            if (s.track == null || s.track.kind === "video") {
                s.replaceTrack(track);
            }
        }
        for (const t of peerConnection.getTransceivers()) {
            if (t.sender.track?.kind == null ||
                t.sender.track.kind === "video") {
                t.direction = "sendrecv";
            }
        }
        if (this.connectionUpdatedEvent != null) {
            this.connectionUpdatedEvent();
        }
    }
    private removeVideoTrack(peerConnection: RTCPeerConnection) {
        const senders = peerConnection.getSenders();
        if (senders.length > 0) {
            this.localVideo.pause();
            for (const s of senders) {
                if (s.track?.kind === "video") {
                    peerConnection.removeTrack(s);
                }
            }
            if (this.connectionUpdatedEvent != null) {
                this.connectionUpdatedEvent();
            }
        }
    }
}