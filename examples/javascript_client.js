
// Set the server and auth key
const SERVER_DOMAIN = "wsrelay.example.com"
const AUTH_KEY = "..." // Do not hardcode auth keys in client-side code that is visible to the user

const ws = new WebSocket(`wss://${SERVER_DOMAIN}/?auth=${encodeURIComponent(AUTH_KEY)}`)

function sendEvent(event_name, data) {
    ws.send(JSON.stringify({event: event_name, data}))
}

// Do this when the connection is successfully established
ws.onopen = (_) => {
    console.log("Connected to server")
};

// Do this when the connection closes
ws.onclose = (event) => {
    console.log(`Connection closed with reason: ${event.reason}`)
};

// Handle receiving events
ws.onmessage = (msg_event) => {
    let message = JSON.parse(msg_event.data);

    switch (message.event) {
        case "log":
            // ...
            break;
        case "new_data":
            // ...
            break;
    }
};

// Send an event
sendEvent("button_click", {
    button_id: 12,
    duration: 204,
});

