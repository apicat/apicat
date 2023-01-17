import { cloneDeep } from 'lodash-es'

/**
 * @name 从一维数组中找到节点的父祖节点
 * @param {Object} item 当前节点
 * @param {Array} arr 全部节点数组
 * @param {Object} options 配置项
 */
export const findParents = (item, arr, options = { id: 'id', parentId: 'parentId', root: 0 }) => {
    let _parents = []
    return (function findParent(item) {
        if (item[options.parentId] === options.root) return _parents
        const parent = arr.find((i) => i[options.id] === item[options.parentId])
        if (parent) {
            _parents.push(parent)
            return findParent(parent)
        } else {
            return _parents
        }
    })(item)
}

/**
 * 将树形数据向下递归为一维数组
 * @param {*} arr 数据源
 * @param {*} childKey  子集key
 */
export const flattenDeep = (arr = [], childKey = 'children') => {
    return arr.reduce((flat, obj) => {
        const item = cloneDeep(obj)
        return flat.concat(item, item[childKey] ? flattenDeep(item[childKey], childKey) : [])
    }, [])
}

let sortVal = 0
export const traverseForMarkSort = (node) => {
    let stack = []
    if (node) {
        stack.push(node)
        while (stack.length) {
            let item = stack.shift()
            item.sort = sortVal
            sortVal++
            let children = item.sub_nodes || []
            for (let i = 0; i < children.length; i++) {
                stack.push(children[i])
            }
        }
    }
}

export const arrayToTree = (array = [], options = { id: 'id', pid: 'pid', children: 'children', rootPidVal: null }) => {
    const res = []
    let itemMap = {}
    ;(array || []).forEach((item) => {
        const { id, parent_id } = item
        itemMap[id] || (itemMap[id] = { sub_nodes: [] })
        itemMap[id] = {
            ...item,
            sub_nodes: itemMap[id]['sub_nodes'],
        }
        const treeItem = itemMap[id]
        if (parent_id === 0) {
            res.push(treeItem)
        } else {
            itemMap[parent_id] || (itemMap[parent_id] = { sub_nodes: [] })
            itemMap[parent_id]['sub_nodes'].push(treeItem)
        }
    })

    return res
}
