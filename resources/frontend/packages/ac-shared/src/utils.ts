export function isFunction(val: unknown): boolean {
    return typeof val === 'function'
}

// 获取路由基础信息
export const getRouteNormalInfo = (route: Array<any>) => route.map((item) => ({ name: item.name, meta: item.meta }))

export const getFirstChar = (str: string) => (str || '').substring(0, 1).toUpperCase()

export const traverseTree = (cb: any, nodes: any, props = { subKey: 'children' }): Array<any> | boolean | void => {
    if (!nodes) return

    const { subKey } = props

    let shouldStop = false

    const result: any = []

    for (let nodeInd = 0; nodeInd < nodes.length; nodeInd++) {
        const node = nodes[nodeInd]

        shouldStop = cb(node, nodes) === false

        result.push(node)

        if (shouldStop) break

        const children: any = node[subKey]

        if (children && children.length) {
            shouldStop = traverseTree(cb, children, props) === false
            if (shouldStop) break
        }
    }

    return !shouldStop ? result : false
}

export function toggleClass(el: any, className: any) {
    el && el.classList && el.classList.toggle(className)
}

export function getAttr(el: any, attr: any) {
    return el && el.getAttribute(attr)
}

export const classNameToArray = (cls = '') => cls.split(' ').filter((item) => !!item.trim())

export const hasClass = (el: Element, cls: string): boolean => {
    if (!el || !cls) return false
    if (cls.includes(' ')) throw new Error('className should not contain space.')
    return el.classList.contains(cls)
}

export const addClass = (el: Element, cls: string) => {
    if (!el || !cls.trim()) return
    el.classList.add(...classNameToArray(cls))
}

export const removeClass = (el: Element, cls: string) => {
    if (!el || !cls.trim()) return
    el.classList.remove(...classNameToArray(cls))
}

export function showOrHide(el: any, isShow: any) {
    if (!el) {
        return
    }
    el.style.display = isShow ? null : 'none'
}

export const timestampFormat = (dateStr: string) => Math.floor((Date.now() - new Date(dateStr).getTime()) / (1000 * 3600 * 24))
