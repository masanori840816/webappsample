export type DataChannel = {
    dataChannel: RTCDataChannel,
    type: "text"|"binary",
    close: () => void
};
export function createTextDataChannel(label: string, id: number, 
    peerConnection: RTCPeerConnection,
    onmessage: (message: string) => void): DataChannel {
    const dc = peerConnection.createDataChannel(label, {
        id,
        negotiated: true,
        ordered: true
    });
    const decoder = new TextDecoder("utf-8");
    dc.onmessage = (ev) => {        
        const message = decoder.decode(new Uint8Array(ev.data));
        onmessage(message);
    };
    dc.onerror = (ev) => {
        console.error(ev);
    }
    return {
        dataChannel: dc,
        type: "text",
        close: () => {
            dc.close();
        }
    }
}