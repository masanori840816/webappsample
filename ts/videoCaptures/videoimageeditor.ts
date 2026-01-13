export function captureVideo(video: HTMLVideoElement, outputImage: HTMLImageElement): boolean {
    const canvas = document.getElementById("video-capture-canvas") as HTMLCanvasElement;
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;

    const context = canvas.getContext("2d");
    if (context != null) {
        // capture current video frame
        context.drawImage(video, 0, 0, canvas.width, canvas.height);
        // create image
        const dataUrl = canvas.toDataURL("image/png");
        outputImage.src = dataUrl;
        outputImage.style.display = "block";
        return true;
    }
    return false;
}
export function savePhoto(callback: ((pohtoData: Uint8Array) => void)) {
    const canvas = document.getElementById("video-capture-canvas") as HTMLCanvasElement;
    const context = canvas.getContext("2d");
    if(context == null) {
        console.error("Failed to get the Canvas Context");
        return;
    }
    // TODO: add drawn lines, images, etc.
    canvas.toBlob(async (b) => {
        if(b == null) {
            console.error("Failed to convert to Blob");
            return;
        }
        const result = new Uint8Array(await b.arrayBuffer());
        console.log(b);
        callback(result);
    });
}