const { Send } = require('../client')
const { toResp } = require('../types')

async function SendSet() {
    const message = toResp(['SET', 'my-key', 'my-value'])
    const expected = "+OK\r\n"

    try {
        const got = await Send(message);

        if (got !== expected) {
            console.error({ expected, got })
        } else {
            console.log(got.replaceAll("\r\n", "\\r\\n"))
        }
    } catch (err) {
        console.error(err)
    }
}

async function SendGet() {
    const message = toResp(['GET', 'my-key'])
    const expected = "$8\r\nmy-value\r\n"

    try {
        const got = await Send(message);

        if (got !== expected) {
            console.error({ expected, got })
        } else {
            console.log(got.replaceAll("\r\n", "\\r\\n"))
        }
    } catch (err) {
        console.error(err)
    }
}

(async () => {
    await SendSet()
    await SendGet()
})()
