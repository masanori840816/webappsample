export function captureVideo(video: HTMLVideoElement, outputImage: HTMLImageElement): boolean {
    const canvas = document.createElement("canvas");
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
