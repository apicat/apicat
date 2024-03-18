import cloneDeep from 'lodash-es/cloneDeep'

interface Options {
  id: string
  parentId: string
  root: number
}

/**
 * @name 从一维数组中找到节点的父祖节点
 * @param item 当前节点
 * @param arr 全部节点数组
 * @param options 配置项
 */
export function findParents(item: any, arr: any[], options: Options = { id: 'id', parentId: 'parentId', root: 0 }): any[] {
  const _parents: any[] = []
  return (function findParent(item: any): any[] {
    if (item[options.parentId] === options.root)
      return _parents
    const parent = arr.find(i => i[options.id] === item[options.parentId])
    if (parent) {
      _parents.push(parent)
      return findParent(parent)
    }
    else {
      return _parents
    }
  })(item)
}

/**
 * 将树形数据向下递归为一维数组
 * @param arr 数据源
 * @param childKey  子集key
 */
export function flattenDeep(arr: any[] = [], childKey: string = 'children'): any[] {
  return arr.reduce((flat, obj) => {
    const item = cloneDeep(obj)
    return flat.concat(item, item[childKey] ? flattenDeep(item[childKey], childKey) : [])
  }, [])
}

let sortVal = 0
export function traverseForMarkSort(node: any): void {
  const stack: any[] = []
  if (node) {
    stack.push(node)
    while (stack.length) {
      const item = stack.shift()
      item.sort = sortVal
      sortVal++
      const children = item.items || []
      for (let i = 0; i < children.length; i++)
        stack.push(children[i])
    }
  }
}

export function arrayToTree(array: any[] = []): any[] {
  const res: any[] = []
  const itemMap: { [key: string]: any } = {};
  (array || []).forEach((item) => {
    const { id, parent_id } = item
    itemMap[id] || (itemMap[id] = { items: [] })
    itemMap[id] = {
      ...item,
      items: itemMap[id].items,
    }
    const treeItem = itemMap[id]
    if (parent_id === 0) {
      res.push(treeItem)
    }
    else {
      itemMap[parent_id] || (itemMap[parent_id] = { items: [] })
      itemMap[parent_id].items.push(treeItem)
    }
  })

  return res
}
