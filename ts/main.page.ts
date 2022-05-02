import { WebsocketMessage } from "./websocket.type";

let ws: WebSocket|null = null;
export function connect(url: string): void {
    ws = new WebSocket(url);
    ws.onopen = () => sendMessage({
        messageType: "text",
        data: "connected",
    });
    ws.onmessage = data => {
        const message = <WebsocketMessage>JSON.parse(data.data);
        switch(message.messageType) {
            case "text":
                if(typeof message.data === "string") {
                    addReceivedMessage(message.data);
                } else {
                    console.error(message.data);                  
                }
                break;
            default:
                console.log(data);
                break;
        }
    };
}
export function send() {
    const messageArea = document.getElementById("input_message") as HTMLTextAreaElement;
    sendMessage({
        messageType: "text",
        data: messageArea.value,
    });
}
export function close() {
    if(ws == null) {
        console.warn("WebSocket was null");
        return;
    }
    ws.close();
    ws = null;
}
function addReceivedMessage(message: string) {
    const receivedMessageArea = document.getElementById("received_text_area") as HTMLElement;
    const child = document.createElement("div");
    child.textContent = message;
    receivedMessageArea.appendChild(child);
}
function sendMessage(message: WebsocketMessage) {
    if (ws == null) {
        console.warn("WebSocket was null");
        return;
    }
    ws.send(JSON.stringify(message));
}
