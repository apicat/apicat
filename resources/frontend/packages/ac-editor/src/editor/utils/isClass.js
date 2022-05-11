export default function isClass(item) {
    if (item.constructor?.toString().substring(0, 5) !== 'class') {
        return false
    }

    return true
}
