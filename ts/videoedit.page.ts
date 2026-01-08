import * as videoCapture from "./videoCaptures/videoimageeditor";

window.VideoPage = {
    capture() {
        const video = document.getElementById("remote_video") as HTMLVideoElement;
        const outputImage = document.getElementById("captured_image") as HTMLImageElement;
        const captureResult = videoCapture.captureVideo(video, outputImage);
        if(captureResult !== true) {
            alert("failed capturing");
        }
        console.log("OK");
    },
    download() {
        console.log("donwloadfile");
    },
}