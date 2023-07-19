class ServerConnection {
    constructor() {
        this.socket = new WebSocket("ws://localhost:8080/process");
    }
    // send message to server
    send(message) {
        this.socket.send(message);
    }
    // receive message from server
    onReceive(callback) {
        this.socket.onmessage = function (event) {
            callback(event.data);
        }
    }
}