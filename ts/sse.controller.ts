import { hasAnyTexts } from "./hasAnyTexts";
import { ClientMessage } from "./webrtc.type";

export class SseController {
    private es: EventSource|null = null;
    private messageReceivedEvent: ((value: string) => void)|null = null;
    private baseUrl = "";

    public constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }
    public connect(userName: string) {

        this.es = new EventSource(`${this.baseUrl}/sse?user=${userName}`);
        
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
    public addEvents(messageReceivedEvent: (value: string) => void) {
        this.messageReceivedEvent = messageReceivedEvent;
    }
    public sendMessage(message: ClientMessage) {
        fetch(`${this.baseUrl}/sse/message`, 
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