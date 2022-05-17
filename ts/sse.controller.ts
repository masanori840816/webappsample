import { hasAnyTexts } from "./hasAnyTexts";
import { ClientMessage } from "./webrtc.type";

export class SseController {
    private es: EventSource|null = null;
    private connectedEvent: (() => void)|null = null;
    private messageReceivedEvent: ((value: string) => void)|null = null;

    public connect(userName: string) {
        this.es = new EventSource(`http://localhost:8080/sse?user=${userName}`);
        
        this.es.onopen = (ev) => {
            console.log(ev);
            if(this.connectedEvent != null) {
                this.connectedEvent();
            }
        };
        this.es.onmessage = (ev) => {
            if(!hasAnyTexts(ev.data) ||
                this.messageReceivedEvent == null) {
                return;
            }
            this.messageReceivedEvent(ev.data);
        };
        this.es.onerror = ev => {
            console.error(ev);        
        };
    }
    public addEvents(connectedEvent: () => void,
        messageReceivedEvent: (value: string) => void) {
        this.connectedEvent = connectedEvent;
        this.messageReceivedEvent = messageReceivedEvent;
    }
    public sendMessage(message: ClientMessage) {
        fetch("http://localhost:8080/sse/message", 
            {
                method: "POST",
                body: JSON.stringify(message),
            })
            .then(res => res.json())
            .then(json => console.log(json))
            .catch(err => console.error(err));
    }
    public close() {
        this.es?.close();
    }
}