const { Send } = require('../client')
const { toResp } = require('../types')

async function SendEcho() {
    const message = toResp(['ECHO', 'hey'])
    const expected = "$3\r\nhey\r\n"

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

(async () => await SendEcho())()
