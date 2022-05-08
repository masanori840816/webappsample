let targetUrl = ""
export function connect(url: string): void {
    targetUrl = url;
    fetch(`${targetUrl}?id=32&name=Hello`, 
    {
        method: "GET",
    })
    .then(res => res.json())
    .then(json => console.log(json))
    .catch(err => console.error(err));
}
export function send() {
    fetch(targetUrl, 
        {
            method: "POST",
            body: JSON.stringify({ messageType: "text", data: "Hello" }),
        })
        .then(res => res.json())
        .then(json => console.log(json))
        .catch(err => console.error(err));
}
export function close() {
    console.log("Close");
    
}