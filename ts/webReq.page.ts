export function getMessage(url: string): void {
    fetch(`${url}?name=nonome3.png`, 
    {
        method: "GET",
    })
    .then(res => res.text())
    .then(txt => console.log(txt))
    .catch(err => console.error(err));
}
export function postMessage(url: string): void {
    fetch(url, 
        {
            method: "POST",
            body: JSON.stringify({ messageType: "text", data: "Hello" }),
        })
        .then(res => res.json())
        .then(json => console.log(json))
        .catch(err => console.error(err));
}