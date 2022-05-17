export class WebRtcController {
    public init() {
        const localVideo = document.getElementById("local_video") as HTMLVideoElement;
        localVideo.addEventListener('canplay', () => {
            const width = 320;
            const height = localVideo.videoHeight / (localVideo.videoWidth/width);          
            localVideo.setAttribute('width', width.toString());
            localVideo.setAttribute('height', height.toString());
          }, false);
        navigator.mediaDevices.getUserMedia({ video: true, audio: true })
          .then(stream => {
              localVideo.srcObject = stream;
              localVideo.play();
          })
          .catch(err => console.error(`An error occurred: ${err}`));
    }
}