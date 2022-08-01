export type ClientMessage = {
    event: "text"|"offer"|"answer"|"candidate"|"update",
    userName: string,
    data: string,
};