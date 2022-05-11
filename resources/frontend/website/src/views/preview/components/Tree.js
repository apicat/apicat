import { DOCUMENT_TYPES } from '@ac/shared'

export class TreeNode {
    constructor(data) {
        const { id, isLeaf, editable, expanded, selected } = data
        this.id = typeof id === 'undefined' ? new Date().valueOf() : id
        this.parent = null
        this.children = null
        this.isLeaf = !!isLeaf
        this.isHover = false
        this.editable = !!editable
        this.expanded = !!expanded
        this.selected = !!selected
        this.dragDisabled = false

        // other params
        for (var k in data) {
            if (k !== 'id' && k !== 'children' && k !== 'isLeaf') {
                this[k] = data[k]
            }
        }
    }

    changeName(name) {
        this.name = name
    }

    addChildren(children, isFirst = false) {
        if (!this.children) {
            this.children = []
        }

        if (Array.isArray(children)) {
            for (let i = 0, len = children.length; i < len; i++) {
                const child = children[i]
                child.parent = this
                child.pid = this.id
            }
            this.children.concat(children)
        } else {
            const child = children
            child.parent = this
            child.pid = this.id
            child.level = this.level + 1

            this.children[!isFirst ? 'push' : 'unshift'](child)
        }

        this.sort(this)
    }

    // remove self
    remove() {
        const parent = this.parent
        const index = parent.findChildIndex(this)
        parent.children.splice(index, 1)
        this.sort(parent)
    }

    // remove child
    _removeChild(child) {
        for (var i = 0, len = this.children.length; i < len; i++) {
            if (this.children[i] === child) {
                this.children.splice(i, 1)
                break
            }
        }
    }

    isTargetChild(target) {
        let parent = target.parent
        while (parent) {
            if (parent === this) {
                return true
            }
            parent = parent.parent
        }
        return false
    }

    moveInto(target) {
        if (this.name === 'root' || this === target) {
            return
        }

        // cannot move ancestor to child
        if (this.isTargetChild(target)) {
            return
        }

        // cannot move to leaf node
        if (target.isLeaf) {
            return
        }

        this.parent._removeChild(this)
        this.parent = target
        this.pid = target.id
        if (!target.children) {
            target.children = []
        }
        target.children.unshift(this)

        this.sort(target)
    }

    findChildIndex(child) {
        var index
        for (let i = 0, len = this.children.length; i < len; i++) {
            if (this.children[i] === child) {
                index = i
                break
            }
        }
        return index
    }

    _canInsert(target) {
        if (this.name === 'root' || this === target) {
            return false
        }
        // cannot insert ancestor to child
        if (this.isTargetChild(target)) {
            return false
        }

        this.parent._removeChild(this)
        this.parent = target.parent
        this.pid = target.parent.id
        return true
    }

    insertBefore(target) {
        if (!this._canInsert(target)) return

        const pos = target.parent.findChildIndex(target)
        target.parent.children.splice(pos, 0, this)
        this.sort(target.parent)
    }

    insertAfter(target) {
        if (!this._canInsert(target)) return

        const pos = target.parent.findChildIndex(target)
        target.parent.children.splice(pos + 1, 0, this)

        this.sort(target.parent)
    }

    insertAfterNoSwap(target) {
        let parent = this.parent
        let pos = parent.findChildIndex(this)
        parent.children.splice(pos + 1, 0, target)
        this.sort(parent)
    }

    sort(sourceTree) {
        // 当前层级
        let level = sourceTree.level
        // 子集层级
        ;(sourceTree.children || []).forEach((node, index) => {
            node.index = index
            node.level = level + 1
            if (node.children && node.children.length) {
                this.sort(node)
            }
        })
    }
}

export class Tree {
    constructor(data) {
        this.root = new TreeNode({ name: 'root', isLeaf: false, id: 0, isRoot: true, level: 0 })
        this.getNodes(this.root, data)
        return this
    }

    getNodes(parent, nodes, parentPath = []) {
        ;(nodes || []).forEach((item, index) => {
            var nodePath = parentPath.concat(index)

            var data = this.getNode(nodePath, item, index)

            var child = new TreeNode(data)

            parent.addChildren(child)

            if (item.sub_nodes && item.sub_nodes.length > 0) {
                this.getNodes(child, item.sub_nodes, nodePath)
            }
        })
    }

    getNode(path, node, index) {
        return {
            index,
            isLeaf: !(node.type === DOCUMENT_TYPES.DIR),
            name: node.title || '',
            ...node,
            _id: node.doc_type + node.id,
            path,
            level: path.length,
        }
    }

    static traverse(cb, nodes = null, parentPath = []) {
        if (!nodes) return

        let shouldStop = false

        const result = []

        for (let nodeInd = 0; nodeInd < nodes.length; nodeInd++) {
            const node = nodes[nodeInd]
            const itemPath = parentPath.concat(nodeInd)

            shouldStop = cb(node, nodes) === false
            result.push(node)

            if (shouldStop) break

            if (node.children) {
                shouldStop = this.traverse(cb, node.children, itemPath) === false
                if (shouldStop) break
            }
        }

        return !shouldStop ? result : false
    }

    static getDepth(tree) {
        if (!tree.children || !tree.children.length) {
            return 1
        }

        let maxDepth = (child) => {
            let deep = 1
            child.forEach((item) => {
                if (item.children) {
                    deep = Math.max(deep, maxDepth(item.children) + 1)
                } else {
                    deep = deep + 1
                }
            })
            return deep
        }

        return maxDepth(tree.children)
    }

    static findNode(id, nodes) {
        let node = null
        Tree.traverse((_node) => {
            if (_node.id === id) {
                node = _node
                return false
            }
        }, nodes)

        return node
    }

    /**
     * returns 1 if path1 > path2
     * returns -1 if path1 < path2
     * returns 0 if path1 == path2
     *
     * examples
     *
     * [1, 2, 3] < [1, 2, 4]
     * [1, 1, 3] < [1, 2, 3]
     * [1, 2, 3] > [1, 2, 0]
     * [1, 2, 3] > [1, 1, 3]
     * [1, 2] < [1, 2, 0]
     *
     */
    comparePaths(path1, path2) {
        for (let i = 0; i < path1.length; i++) {
            if (path2[i] == void 0) return 1
            if (path1[i] > path2[i]) return 1
            if (path1[i] < path2[i]) return -1
        }
        return path2[path1.length] == void 0 ? 0 : -1
    }

    getNextNode(path, filter = null) {
        let resultNode = null

        Tree.traverse((node) => {
            if (this.comparePaths(node.path, path) < 1) return

            if (!filter || filter(node)) {
                resultNode = node
                return false // stop traverse
            }
        }, this.root.children)

        return resultNode
    }

    getPrevNode(path, filter) {
        let prevNodes = []

        Tree.traverse((node) => {
            if (this.comparePaths(node.path, path) >= 0) {
                return false
            }
            prevNodes.push(node)
        }, this.root.children)

        let i = prevNodes.length
        while (i--) {
            const node = prevNodes[i]
            if (!filter || filter(node)) return node
        }

        return null
    }
}

export default Tree
