import * as urlParam from "./urlParamGetter";
import { ClientName, ClientNames } from "./webrtc.type";
type RemoteTrack = {
    id: string,
    kind: "video"|"audio",
    element: HTMLElement,
};
type ConnectedClient = {
    name: ClientName,
    element: HTMLElement,
};
export class MainView {
    private receivedTextArea: HTMLElement;
    private receivedDataChannelArea: HTMLElement;
    private localVideoUsed: HTMLInputElement;
    private remoteTrackArea: HTMLElement;
    private tracks = new Array<RemoteTrack>();
    private clientArea: HTMLElement;
    private connectedClients: ConnectedClient[];
    private forceVideoCodecCheck: HTMLInputElement;
    private preferVideoCodecSelector: HTMLSelectElement;
    public getForceVideoCodec(): boolean {
        return this.forceVideoCodecCheck.checked;
    }
    public getPreferredVideoCodec(): string {
        return this.preferVideoCodecSelector.value;
    }
    public constructor() {
        this.receivedTextArea = document.getElementById("received_text_area") as HTMLElement;
        this.receivedDataChannelArea = document.getElementById("received_datachannel_area") as HTMLElement;
        this.localVideoUsed = document.getElementById("local_video_usage") as HTMLInputElement;
        this.localVideoUsed.checked = urlParam.getBoolParam("video");
        this.remoteTrackArea = document.getElementById("remote_track_area") as HTMLElement;
        this.clientArea = document.getElementById("client_names") as HTMLElement;
        this.connectedClients = new Array<ConnectedClient>();
        this.forceVideoCodecCheck = document.getElementById("force_video_codec_check") as HTMLInputElement;
        this.preferVideoCodecSelector = document.getElementById("video_codec_selector") as HTMLSelectElement;
        this.addVideoCodecNames();
    }
    public addEvents(videoUsageChanged: (used: boolean) => void): void {
        this.localVideoUsed.onchange = () =>
            videoUsageChanged(this.localVideoUsed.checked);
    }
    public addReceivedText(value: { user: string, message: string}): void {
        const newText = document.createElement("div");
        newText.textContent = `User: ${value.user} Message: ${value.message}`;
        this.receivedTextArea.appendChild(newText);
    }
    public addReceivedDataChannelValue(value: string|Uint8Array): void {
        if(typeof value === "string") {
            const newText = document.createElement("div");
            newText.textContent = value;
            this.receivedDataChannelArea.appendChild(newText);
            return;
        }
        console.warn("Not implemented");
    }
    public checkLocalVideoUsed(): boolean {
        return this.localVideoUsed.checked;
    }
    public addRemoteTrack(stream: MediaStream, kind: "video"|"audio", id?: string): void {
        if(this.tracks.some(t => t.id === stream.id)) {
            if(kind === "audio") {
                return;
            }
            this.removeRemoteTrack(stream.id, "audio");
        }
        const remoteTrack = document.createElement(kind);
        remoteTrack.srcObject = stream;
        remoteTrack.autoplay = true;
        remoteTrack.controls = false;
        this.remoteTrackArea.appendChild(remoteTrack);
        this.tracks.push({
            id: (id == null)? stream.id: id,
            kind,
            element: remoteTrack,
        });        
    }
    public removeRemoteTrack(id: string, kind: "video"|"audio"): void {
        const targets = this.tracks.filter(t => t.id === id);
        if(targets.length <= 0) {
            return;
        }
        if(kind === "video") {
            const audioTrack = this.getAudioTrack(targets[0]?.element);
            if(audioTrack != null) {
                this.addRemoteTrack(new MediaStream([audioTrack]), "audio", id);
            }
        }
        for(const t of targets) {
            this.remoteTrackArea.removeChild(t.element);
        }
        const newTracks = new Array<RemoteTrack>();
        for(const t of this.tracks.filter(t => t.id !== id || (t.id === id && t.kind !== kind))) {
            newTracks.push(t);
        }
        this.tracks = newTracks;
    }
    public updateClientNames(names: ClientNames): void {
        if(names == null) {
            console.warn("updateClientNames were null");
            return;
        }
        const newClients = new Array<ConnectedClient>();
        for(const c of this.connectedClients) {
            const clientName = c.name.name;
            if(names.names.some(n => n.name === clientName)) {
                newClients.push(c);
            } else {
                this.clientArea.removeChild(c.element);
            }
        }
        for(const n of names.names) {
            const clientName = n;
            if(this.connectedClients.some(c => c.name.name === clientName.name) === false) {
                const newElement = document.createElement("div");
                newElement.textContent = clientName.name;
                this.clientArea.appendChild(newElement);
                this.connectedClients.push({
                    name: clientName,
                    element: newElement,
                });
            }
        }
    }
    private getAudioTrack(target: HTMLElement|null|undefined): MediaStreamTrack|null {
        if(target == null ||
            !(target instanceof HTMLVideoElement)){
            return null;
        }
        if(target.srcObject == null ||
            !("getAudioTracks" in target.srcObject) ||
            (typeof target.srcObject.getAudioTracks !== "function")) {
            return null;
        }
        const tracks = target.srcObject.getAudioTracks();
        if(tracks.length <= 0 ||
            tracks[0] == null) {
            return null;
        }
        return tracks[0];
    }
    private addVideoCodecNames() {
        const codecs = RTCRtpSender.getCapabilities("video")?.codecs;
        if(codecs == null) {
            return;
        }
        const addedMimeTypes: string[] = [];
        const option = document.createElement("option");
        option.value = "";
        option.text = "";
        this.preferVideoCodecSelector.appendChild(option);
        for(const c of codecs) {
            if(addedMimeTypes.some(t => c.mimeType === t)) {
                continue;
            }
            addedMimeTypes.push(c.mimeType);
            const option = document.createElement("option");
            option.value = c.mimeType;
            option.text = c.mimeType;
            this.preferVideoCodecSelector.appendChild(option);
        }
    }
}