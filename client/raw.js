const { Send } = require('./client');
const { toResp } = require('./types');

cmd = process.argv[2].match(/(?:[^\s"]+|"[^"]*")+/g)

for (let i = 0; i < cmd.length; i++) {
    if (cmd[i].startsWith('"') && cmd[i].endsWith('"')) {
        cmd[i] = cmd[i].replace(/^"(.*)"$/, '$1')
    }
}

const formattedCmd = toResp(cmd);

(async () => {
    try {
        rawresponse = await Send(formattedCmd)

        console.log(rawresponse.replaceAll('\r\n', '\\r\\n'));
    } catch (e) {
        console.error(e);
    }
})()
