export type TextDataChannel = {
    dataChannel: RTCDataChannel,
    type: "text"|"binary",
    close: () => void
};
export function createTextDataChannel(label: string, id: number, 
    peerConnection: RTCPeerConnection): TextDataChannel {
    const dc = peerConnection.createDataChannel(label, {
        id,
        negotiated: true,
        ordered: true
    });
    dc.onopen = (ev) => {
        console.log("OnOpen");
        console.log(ev);

        dc.send("hellllloooo");
    }
    dc.onmessage = (ev) => {
        console.log("OnMessage");
        console.log(ev);
    };
    dc.onerror = (ev) => {
        console.log("OnError");
        console.log(ev);
    }
    return {
        dataChannel: dc,
        type: "text",
        close: () => {
            dc.close();
        }
    }
}