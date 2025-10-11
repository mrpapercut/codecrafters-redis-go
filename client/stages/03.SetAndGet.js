const { Send } = require('../client');

async function SendSet() {
    const cmd = "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$8\r\nmy-value\r\n";

    try {
        const res = await Send(cmd);

        console.log(res);
    } catch (err) {
        console.error(err)
    }
}

async function SendGet() {
    const cmd = "*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n";

    try {
        const res = await Send(cmd);

        console.log(res);
    } catch (err) {
        console.error(err)
    }
}

(async () => {
    await SendSet()
    await SendGet()
})()
