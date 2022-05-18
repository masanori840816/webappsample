export class MainView {
    private receivedTextArea: HTMLElement;
    public constructor() {
        this.receivedTextArea = document.getElementById("received_text_area") as HTMLElement;
    }
    public addReceivedText(value: { user: string, message: string}): void {
        const newText = document.createElement("div");
        newText.textContent = `User: ${value.user} Message: ${value.message}`;
        this.receivedTextArea.appendChild(newText);
    }
}