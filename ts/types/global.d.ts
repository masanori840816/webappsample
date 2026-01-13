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
    sendTextDataChannel: () => void,
    capture: () => void,
}
export interface VideoEditPageApi {
    init: () => void,
    capture: () => void,
    download: () => void,
}