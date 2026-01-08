import * as videoCapture from "./videoCaptures/videoimageeditor";

let captured = false;
window.VideoPage = {
    capture() {
        const video = document.getElementById("remote_video") as HTMLVideoElement;
        const outputImage = document.getElementById("captured_image") as HTMLImageElement;
        const captureResult = videoCapture.captureVideo(video, outputImage);
        if(captureResult !== true) {
            alert("failed capturing");
            return;
        }
        captured = true;
        console.log("OK");
    },
    download() {
        console.log("donwloadfile");
    },
    mark() {
        if(captured !== true) {
            return;
        }
        console.log("mark");
    }
}