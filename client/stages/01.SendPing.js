const { Send } = require('../client');

async function SendPing() {
    const cmd = "+PING\r\n"

    try {
        const res = await Send(cmd);

        console.log(res);
    } catch (err) {
        console.error(err)
    }
}

(async () => await SendPing())()
