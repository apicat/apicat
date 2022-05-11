export default function objectIncludes(object1, object2) {
    const keys = Object.keys(object2);
    if (!keys.length) {
        return true;
    }
    return !!keys
        .filter(key => object2[key] === object1[key])
        .length;
}
