import { hasAnyTexts } from "./hasAnyTexts";
import { MainView } from "./main.view";
import { SseController } from "./sse.controller";
import { removeVideoCodec } from "./videoCodecs/videoCodecRemover";
import { WebRtcController } from "./webrtc.controller";
import { ClientMessage, ICEServer } from "./webrtc.type";
import * as videoCapture from "./videoCaptures/videoimageeditor";
import { toBase64 } from "./videoCaptures/blobConverter";

let sse: SseController;
let webrtc: WebRtcController;
let view: MainView;
let userName = ""
let iceServer: ICEServer;

let captured = false;
// TODO: get from the page
const targetSample = "sample";
window.Page = {
    connect(): void {
        const userNameInput = document.getElementById("user_name") as HTMLInputElement;
        userName = userNameInput.value;
        webrtc.connect(iceServer);
        sse.connect(userName, targetSample);
    },
    send() {
        if(!hasAnyTexts(userName)) {
            return;
        }
        const messageInput = document.getElementById("input_message") as HTMLTextAreaElement;
        sse.sendMessage({ event: "text", userName, groupName: targetSample, data: messageInput.value });
    },
    close() {
        userName = "";
        webrtc.close();
        sse.close();
    },
    init(url: string, iceServerJSON: string) {
        iceServer = JSON.parse(iceServerJSON);       
        sse = new SseController(url);
        sse.addEvents((value) => handleReceivedMessage(value));
        
        view = new MainView();
        webrtc = new WebRtcController();
        webrtc.addEvents((message) => sendAnswer(message),
            (message) => sendCandidate(message),
            (message) => view.addReceivedDataChannelValue(message),
            () => updateConnection(),
            (stream, kind) => view.addRemoteTrack(stream, kind),
            (id, kind) => view.removeRemoteTrack(id, kind));
        webrtc.init(view.checkLocalVideoUsed());
        view.addEvents((used) => webrtc.switchLocalVideoUsage(used));
    },
    sendTextDataChannel() {
        const messageInput = document.getElementById("input_message") as HTMLTextAreaElement;
        webrtc.sendTextDataChannel(messageInput.value);
    },
    capture() {
        const video = document.getElementById("remote_video") as HTMLVideoElement;
        const outputImage = document.getElementById("captured_image") as HTMLImageElement;
        const captureResult = videoCapture.captureVideo(video, outputImage);
        if(captureResult !== true) {
            alert("failed capturing");
            return;
        }
        captured = true;
        console.log("OK");
        videoCapture.savePhoto((photoData) => {
            if(captured !== true) {
                return;
            }
            console.log("video capture send");
            console.log(photoData);
            
            sse.sendMessage({
                event: "photo",
                userName, 
                groupName: targetSample, 
                data: toBase64(photoData)
            });
        });
    },

};
function handleReceivedMessage(value: string) {
    const message = JSON.parse(value);
    if(!checkIsClientMessage(message)) {
        console.error(`Invalid message type ${value}`);
        return;
    }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    let offerData: any;
    let preferredMimeType: string|null = null;
    switch(message.event) {
        case "text":
            view.addReceivedText({ user: message.userName, message: message.data });
            break;
        case "offer":
            offerData = JSON.parse(message.data);
            if(view.getForceVideoCodec()) {
                offerData.sdp = removeVideoCodec(offerData.sdp, view.getPreferredVideoCodec());
            } else {
                preferredMimeType = view.getPreferredVideoCodec();
            }       
            if(webrtc.handleOffer(offerData, preferredMimeType) !== true) {
                close();
            }
            break;
        case "candidate":
            webrtc.handleCandidate(JSON.parse(message.data));
            break;
        case "clientName":
            view.updateClientNames(JSON.parse(message.data));
            break;
        case "heartbeat":
            // Do nothing
            break;
        case "photo":
            console.log(`recevie photo ${message.data}`);
            break;
        default:
            console.error(`Invalid message type ${value}`);            
            break;
    }
}
function sendAnswer(data: RTCSessionDescriptionInit) {
    if(!hasAnyTexts(userName)) {
        return;
    }
    sse.sendMessage({userName, event: "answer", groupName: targetSample, data: JSON.stringify(data)});
}
function sendCandidate(data: RTCIceCandidate) {
    if(!hasAnyTexts(userName)) {
        return;
    }
    sse.sendMessage({userName, event: "candidate", groupName: targetSample, data: JSON.stringify(data)});
}
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function checkIsClientMessage(value: any): value is ClientMessage {
    if(value == null) {
        return false;
    }
    if(("event" in value &&
        typeof value["event"] === "string") === false) {
        return false;
    }
    if(("data" in value &&
        typeof value["data"] === "string") === false) {
        return false;
    }
    return true;
}
function updateConnection() {
    if(!hasAnyTexts(userName)) {
        return;
    }
    sse.sendMessage({userName, event: "update", groupName: targetSample, data: "{}"});
}