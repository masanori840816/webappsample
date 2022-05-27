import { hasAnyTexts } from "./hasAnyTexts";
import { MainView } from "./main.view";
import { SseController } from "./sse.controller";
import { WebRtcController } from "./webrtc.controller";
import { ClientMessage } from "./webrtc.type";

let sse: SseController;
let webrtc: WebRtcController;
let view: MainView;
let userName = ""
export function connect(): void {
    const userNameInput = document.getElementById("user_name") as HTMLInputElement;
    userName = userNameInput.value;
    sse.connect(userName);
}
export function send() {
    if(!hasAnyTexts(userName)) {
        return;
    }
    const messageInput = document.getElementById("input_message") as HTMLTextAreaElement;
    sse.sendMessage({ event: "text", userName, data: messageInput.value });
}
export function close() {
    userName = "";
    sse.close();
}
export function init(url: string) {
    sse = new SseController(url);
    sse.addEvents((value) => handleReceivedMessage(value));
    webrtc = new WebRtcController();
    webrtc.addEvents((message) => sendAnswer(message),
        (message) => sendCandidate(message),
        (message) => view.addReceivedDataChannelValue(message));
    webrtc.init();
    view = new MainView();
}
export function sendTextDataChannel() {
    const messageInput = document.getElementById("input_message") as HTMLTextAreaElement;
    webrtc.sendTextDataChannel(messageInput.value);
}
function handleReceivedMessage(value: string) {
    const message = JSON.parse(value);
    if(!checkIsClientMessage(message)) {
        console.error(`Invalid message type ${value}`);
        return;
    }
    switch(message.event) {
        case "text":
            view.addReceivedText({ user: message.userName, message: message.data });
            break;
        case "offer":
            webrtc.handleOffer(JSON.parse(message.data));
            break;
        case "candidate":
            webrtc.handleCandidate(JSON.parse(message.data));
            break;
    }
}
function sendAnswer(data: RTCSessionDescriptionInit) {
    if(!hasAnyTexts(userName)) {
        return;
    }
    sse.sendMessage({userName, event: "answer", data: JSON.stringify(data)});
}
function sendCandidate(data: RTCIceCandidate) {
    if(!hasAnyTexts(userName)) {
        return;
    }
    sse.sendMessage({userName, event: "candidate", data: JSON.stringify(data)});
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
    if(value.event !== "text" &&
        value.event !== "offer" &&
        value.event !== "answer" &&
        value.event !== "candidate") {
        return false;
    }
    if(("data" in value &&
        typeof value["data"] === "string") === false) {
        return false;
    }
    return true;
}
