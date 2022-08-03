import * as urlParam from "./urlParamGetter";
type RemoteTrack = {
    id: string,
    kind: "video"|"audio",
    element: HTMLElement,
};
export class MainView {
    private receivedTextArea: HTMLElement;
    private receivedDataChannelArea: HTMLElement;
    private localVideoUsed: HTMLInputElement;
    private remoteTrackArea: HTMLElement;
    private tracks = new Array<RemoteTrack>();
    public constructor() {
        this.receivedTextArea = document.getElementById("received_text_area") as HTMLElement;
        this.receivedDataChannelArea = document.getElementById("received_datachannel_area") as HTMLElement;
        this.localVideoUsed = document.getElementById("local_video_usage") as HTMLInputElement;
        this.localVideoUsed.checked = urlParam.getBoolParam("video");
        this.remoteTrackArea = document.getElementById("remote_track_area") as HTMLElement;
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
            console.log(this.tracks);
            
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
}