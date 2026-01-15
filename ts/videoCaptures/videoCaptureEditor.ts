const MarkerSize = 60;
export class VideoCaptureEditor {
    private canvas: HTMLCanvasElement;
    private markingArea: HTMLElement;
    private markerImg: HTMLImageElement;
    private markers: { x: number, y: number }[] = [];

    public constructor() {
        this.canvas = document.createElement("canvas");
        this.markingArea = document.getElementById("marking-capture-area") as HTMLElement;
        this.markerImg = new Image();
        this.markerImg.src = 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(`
            <svg viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                <circle cx="10" cy="10" r="8" fill="none" stroke="red" stroke-width="2" />
            </svg>
        `);
    }
    public capture(video: HTMLVideoElement, outputImage: HTMLImageElement): boolean {

        for(const c of this.markingArea.children)
        {
            this.markingArea.removeChild(c);
        }

        this.canvas.width = video.videoWidth;
        this.canvas.height = video.videoHeight;
        const context = this.canvas.getContext("2d");
        if (context != null) {
            // clear last images
            context.clearRect(0, 0, this.canvas.width, this.canvas.height);
            this.markers = [];
            // capture current video frame
            context.drawImage(video, 0, 0, this.canvas.width, this.canvas.height);
            // create image
            const dataUrl = this.canvas.toDataURL("image/png");
            outputImage.src = dataUrl;
            outputImage.style.display = "block";
            return true;
        }
        return false;
    }
    public savePhoto(outputImage: HTMLImageElement, callback: ((pohtoData: Uint8Array) => void)) {
        const context = this.canvas.getContext("2d");
        if(context == null) {
            console.error("Failed to get the Canvas Context");
            return;
        }
        this.canvas.width = outputImage.naturalWidth;
        this.canvas.height = outputImage.naturalHeight;
        const scaleX = this.canvas.width / outputImage.clientWidth;
        const scaleY = this.canvas.height / outputImage.clientHeight;

        context.drawImage(outputImage, 0, 0);
        const markerSizeOnCanvas = MarkerSize * scaleX;
        this.markers.forEach(m => {
            context.drawImage(
                this.markerImg, 
                (m.x * scaleX) - (markerSizeOnCanvas / 2), 
                (m.y * scaleY) - (markerSizeOnCanvas / 2), 
                markerSizeOnCanvas, 
                markerSizeOnCanvas
            );
        });

        const link = document.createElement('a');
        link.download = 'marked_image.png';
        link.href = this.canvas.toDataURL('image/png');
        link.click();

        // TODO: add drawn lines, images, etc.
        this.canvas.toBlob(async (b) => {
            if(b == null) {
                console.error("Failed to convert to Blob");
                return;
            }
            const result = new Uint8Array(await b.arrayBuffer());
            callback(result);
        });
    }
    public mark(ev: MouseEvent, outputImage: HTMLImageElement) {
        const rect = outputImage.getBoundingClientRect();        
        const x = ev.clientX - rect.left;
        const y = ev.clientY - rect.top;
        this.markers.push({ x, y });
        
        const icon = document.createElement('img');
        icon.src = this.markerImg.src;
        icon.className = 'marker';
        const leftPos = outputImage.offsetLeft + x - (MarkerSize / 2);
        const topPos = outputImage.offsetTop + y - (MarkerSize / 2);
        Object.assign(icon.style, {
            position: 'absolute',
            width: `${MarkerSize}px`,
            height: `${MarkerSize}px`,
            left: `${leftPos}px`,
            top: `${topPos}px`,
            pointerEvents: 'none'
        });
        this.markingArea.appendChild(icon);
    }
}
