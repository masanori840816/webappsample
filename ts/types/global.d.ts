declare global {
    interface Window {
        Page: MainPageApi,
        VideoPage: VideoEditPageApi
    }
}
export interface MainPageApi {
    connect: () => void,
    send: () => void,
    close: () => void,
    init: (url: string, iceServerJSON: string) => void,
    sendTextDataChannel: () => void
}
export interface VideoEditPageApi {
    capture: () => void,
    download: () => void,
}