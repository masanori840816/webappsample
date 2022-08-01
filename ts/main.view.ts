export class MainView {
    private receivedTextArea: HTMLElement;
    private receivedDataChannelArea: HTMLElement;
    private localVideoUsed: HTMLInputElement;
    public constructor() {
        this.receivedTextArea = document.getElementById("received_text_area") as HTMLElement;
        this.receivedDataChannelArea = document.getElementById("received_datachannel_area") as HTMLElement;
        this.localVideoUsed = document.getElementById("local_video_usage") as HTMLInputElement;
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
}