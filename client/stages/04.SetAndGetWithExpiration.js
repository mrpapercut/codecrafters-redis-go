const { Send } = require('../client');

async function SendSet() {
    const cmd = "*5\r\n$3\r\nSET\r\n$6\r\nmy-key\r\n$8\r\nmy-value\r\n$2\r\nEX\r\n$1\r\n1\r\n";

    try {
        const res = await Send(cmd);

        console.log(res);
    } catch (err) {
        console.error(err)
    }
}

async function SendGet() {
    const cmd = "*2\r\n$3\r\nGET\r\n$6\r\nmy-key\r\n";

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
    await (new Promise((res, rej) => setTimeout(res, 500)))
    await SendGet()
    await (new Promise((res, rej) => setTimeout(res, 1000)))
    await SendGet()
})()
