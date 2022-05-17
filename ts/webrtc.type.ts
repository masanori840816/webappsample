export type ClientMessage = {
    event: "text"|"offer"|"answer"|"candidate",
    userName: string,
    data: string,
};
export type CandidateMessage = {
    event: "candidate",
    data: RTCIceCandidate,
};
export type VideoOfferMessage = {
    event: "offer",
    data: RTCSessionDescription|null,
};
export type VideoAnswerMessage = {
    event: "answer",
    data: RTCSessionDescriptionInit|null,
};