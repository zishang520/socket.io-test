/**
 * Module dependencies.
 */

// helper

function $(id) {
    return document.getElementById(id);
}

// chart
for (var i = 0; i < 1; i++) {
    // socket
    const socket = new eio({
        // forceBase64: true,
        transports: ["polling", "websocket", ],
    });
    let last;

    function send() {
        last = new Date();
        socket.send("ping");
        $("transport").innerHTML = socket.transport.name;
    }

    socket.on("open", () => {
        send();
    });

    socket.on("close", () => {
        $("transport").innerHTML = "(disconnected)";
        console.log("已经关闭连接");
    });

    socket.on("message", (m) => {
        console.log(m);
        const latency = new Date() - last;
        $("latency").innerHTML = latency + "ms";
        setTimeout(send, 100);
    });
}
