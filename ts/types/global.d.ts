declare global {
    interface Window {
        Page: MainPageApi,
    }
}
export interface MainPageApi {
    connect: () => void,
    send: () => void,
    close: () => void,
    init: (url: string, iceServerText: string) => void,
    sendTextDataChannel: () => void
}