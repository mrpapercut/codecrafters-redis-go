const net = require('net');

const port = 6379;
const host = '127.0.0.1';

async function Send(command) {
    return new Promise((res, rej) => {
        const socket = new net.Socket();

        socket.connect(port, host);

        socket.on('connect', () => {
            console.log(`Connected to ${host}:${port}`);

            socket.write(command)
        });

        socket.on('data', data => {
            res(data.toString().replaceAll("\r\n", "\\r\\n"));

            socket.destroy();
        });

        socket.on('error', err => {
            console.error('error on socket:', err);
            rej();
        });
    })
}

module.exports = { Send };
