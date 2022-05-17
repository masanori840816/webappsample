import { hasAnyTexts } from "./hasAnyTexts";
import { SseController } from "./sse.controller";
import { WebRtcController } from "./webrtc.controller";
import { CandidateMessage, ClientMessage, VideoAnswerMessage } from "./webrtc.type";

let sse: SseController;
let webrtc: WebRtcController;
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
function init() {
    sse = new SseController();
    sse.addEvents(() => webrtc.connect(),
        (value) => handleReceivedMessage(value));
    webrtc = new WebRtcController();
    webrtc.addEvents((message) => sendAnswer(message),
        (message) => sendCandidate(message));
    webrtc.init();
}
function handleReceivedMessage(value: string) {
    if(!checkIsClientMessage(value)) {
        console.error(`Invalid message type ${value}`);
        return;
    }
    switch(value.event) {
        case "text":
            console.log("Text");
            console.log(value.data);
            break;
        case "offer":
            console.log("Offer");
            webrtc.handleOffer({event: "offer", data: JSON.parse(value.data)});
            break;
        case "answer":
            console.log("Answer");
            break;
        case "candidate":
            console.log("Candidate");
            webrtc.handleCandidate({ event: "candidate", data: JSON.parse(value.data) });
            break;
    }
    console.log(value);
    
}
function sendAnswer(message: VideoAnswerMessage) {
    if(!hasAnyTexts(userName)) {
        return;
    }
    sse.sendMessage({userName, event: message.event, data: JSON.stringify(message.data)});
}
function sendCandidate(message: CandidateMessage) {
    if(!hasAnyTexts(userName)) {
        return;
    }
    sse.sendMessage({userName, event: message.event, data: JSON.stringify(message.data)});
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
init();