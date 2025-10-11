const { Send } = require('../client')
const { toResp } = require('../types')

async function SendRPush() {
    const message = toResp(['RPUSH', 'list_key', 'element'])
    const expected = ":1\r\n"

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

(async () => await SendRPush())()
