export type ClientMessage = {
    event: "text"|"offer"|"answer"|"candidate"|"update"|"clientName",
    userName: string,
    data: string,
};
export type ClientName = {
	name: string,
}
export type ClientNames = {
	names: ClientName[]
}