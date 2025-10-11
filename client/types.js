function toBulkstring(str) {
    return `$${str.length}\r\n${str}\r\n`
}

function toArray(arr) {
    let str = `*${arr.length}\r\n`

    for (let i = 0; i < arr.length; i++) {
        str += toResp(arr[i])
    }

    return str
}

function toResp(msg) {
    if (typeof msg === 'string') {
        return toBulkstring(msg)
    } else if (typeof msg === 'object' && Array.isArray(msg)) {
        return toArray(msg)
    }
}

module.exports = {
    toResp
}
