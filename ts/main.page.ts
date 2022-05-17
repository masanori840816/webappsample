import { hasAnyTexts } from "./hasAnyTexts";
import { SseController } from "./sse.controller";
import { WebRtcController } from "./webrtc.controller";

let sse: SseController;
let webrtc: WebRtcController;
let userName = ""
export function connect(): void {
    const userNameInput = document.getElementById("user_name") as HTMLInputElement;
    userName = userNameInput.value;
    //const receivedArea = document.getElementById("received_text_area") as HTMLElement;
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
    webrtc = new WebRtcController();
    sse.addEvents(() => webrtc.init(),
        (value) => handleReceivedMessage(value));
}
function handleReceivedMessage(value: string) {
    console.log(value);
    
}
init();