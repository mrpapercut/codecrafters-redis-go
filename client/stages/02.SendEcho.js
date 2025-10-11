const { Send } = require('../client');

async function SendEcho() {
    const cmd = "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n";

    try {
        const res = await Send(cmd);

        console.log(res);
    } catch (err) {
        console.error(err)
    }
}

(async () => await SendEcho())()
