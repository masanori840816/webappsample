export type ClientMessage = {
    event: "text"|"offer"|"answer"|"candidate",
    userName: string,
    data: string,
};