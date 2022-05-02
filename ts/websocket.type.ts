export type WebsocketMessage = {
    messageType: "text"
    data: string|Blob|ArrayBuffer,
};