export type ClientMessage = {
    event: "text"|"offer"|"answer"|"candidate"|"update"|"clientName"|"heartbeat",
    userName: string,
    target: string,
    data: string,
};
export type ClientName = {
	name: string,
}
export type ClientNames = {
	names: ClientName[]
}
export type ICEServer = {
    urls: string,
    username: string,
    credential: string,
}