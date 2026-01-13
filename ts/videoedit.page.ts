import { VideoCaptureEditor } from "./videoCaptures/videoCaptureEditor";


let captured = false;
let outputImage: HTMLImageElement;
let videoCapture: VideoCaptureEditor;
window.VideoPage = {
    init() {
        videoCapture = new VideoCaptureEditor();
        outputImage = document.getElementById("captured_image") as HTMLImageElement;
        outputImage.addEventListener("click", (ev) => {
            console.log(ev);
            if(captured !== true) {
                return;
            }
            videoCapture.mark(ev, document.getElementById("captured_image") as HTMLImageElement);
            console.log("mark");
        });
    },
    capture() {
        const video = document.getElementById("remote_video") as HTMLVideoElement;
        const captureResult = videoCapture.capture(video, outputImage);
        if(captureResult !== true) {
            alert("failed capturing");
            return;
        }
        captured = true;
    },
    download() {
        videoCapture.savePhoto(outputImage, (photoData) => console.log(photoData));
    }
}