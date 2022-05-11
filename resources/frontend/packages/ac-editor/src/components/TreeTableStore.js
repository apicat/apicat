import shortid from "shortid";

export class TreeNode {
  constructor(data = {}) {
    //节点所需数据
    this.id = `node_${typeof data.id === "undefined" ? shortid() : data.id}`;
    this.children = null;

    // 是否展开
    this.isExpanded = false;

    // 原数据
    this.data = data;

    // 父节点信息
    this.parent = null;

    // 节点位置信息
    this.path = [];
    this.pathStr = JSON.stringify([]);

    // 层级
    this.level = 0;
  }

  addChildren(children, isInit = false) {
    if (!this.children) {
      this.children = [];
    }

    if (Array.isArray(children)) {
      for (let i = 0, len = this.children.length; i < len; i++) {
        const child = children[i];
        child.parent = this;
        child.level = this.level + 1;
      }
      this.children.concat(children);
    } else {
      const child = children;
      child.parent = this;
      child.level = this.level + 1;
      this.children.push(child);
    }

    !isInit && this.sort();

    this.isExpanded = true;
  }

  remove() {
    if (!this.parent) {
      return null;
    }

    const parent = this.parent;
    const index = parent.findChildIndex(this);
    let result = parent.children.splice(index, 1);

    this.sort();

    return result.length ? result[0] : null;
  }

  removeChild(child) {
    for (let i = 0, len = this.children.length; i < len; i++) {
      if (this.children[i] === child) {
        this.children.splice(i, 1);
        break;
      }
    }
  }

  findChildIndex(child) {
    let index;
    for (let i = 0, len = this.children.length; i < len; i++) {
      if (this.children[i] === child) {
        index = i;
        break;
      }
    }
    return index;
  }

  isTargetChild(target) {
    let parent = target.parent;
    while (parent) {
      if (parent === this) {
        return true;
      }
      parent = parent.parent;
    }
    return false;
  }

  canInsert(target) {
    if (this.isRoot || this === target) {
      return false;
    }

    if (this.isTargetChild(target)) {
      return false;
    }

    this.parent.removeChild(this);
    this.parent = target.parent;
    return true;
  }

  moveInto(target) {
    if (this.isRoot || this === target) {
      return;
    }

    // cannot move ancestor to child
    if (this.isTargetChild(target)) {
      return;
    }

    // cannot move to leaf node
    if (target.isLeaf) {
      return;
    }

    this.parent.removeChild(this);
    this.parent = target;

    if (!target.children) {
      target.children = [];
    }
    target.children.unshift(this);
  }

  insertAtPos(pos, target) {
    if (!this.children) {
      this.children = [];
    }

    this.children.splice(pos, 0, target);
    this.sort();
  }

  insertBefore(target) {
    if (!target.parent) return;
    if (!this.canInsert(target)) return;

    const pos = target.parent.findChildIndex(target);
    target.parent.children.splice(pos, 0, this);
  }

  insertAfter(target) {
    if (!target.parent) return;
    if (!this.canInsert(target)) return;

    const pos = target.parent.findChildIndex(target);
    target.parent.children.splice(pos + 1, 0, this);
  }

  // 是否深度遍历 ？？？
  sort() {
    if (!this.children || !this.children.length) {
      return;
    }

    function dfs(node, data, parentPath = []) {
      for (let i = 0, len = data.length; i < len; i++) {
        let nodePath = parentPath.concat(i);

        let child = data[i];

        // 重新排序字段
        child.path = nodePath;
        child.pathStr = JSON.stringify(nodePath);
        child.level = node.level + 1;

        let children = child.children;
        if (children && children.length > 0) {
          dfs(child, children, nodePath);
        }
      }
    }

    dfs(this, this.children, this.path);
  }

  resort(sourceTree, parentPath = []) {
    // 当前层级
    let level = sourceTree.level;
    // 子集层级
    (sourceTree.children || []).forEach((node, index) => {
      let nodePath = parentPath.concat(index);

      node.path = nodePath;
      node.pathStr = JSON.stringify(nodePath);
      node.level = level + 1;

      if (node.children && node.children.length) {
        this.resort(node, nodePath);
      }
    });
  }
}

export default class TreeTableStore {
  get defaultOption() {
    return {
      childrenKey: "children",
    };
  }

  constructor(data, option) {
    this.option = { ...this.defaultOption, ...option };

    this.root = new TreeNode();
    this.root.isRoot = true;

    this.initTreeNode(this.root, data || []);

    this.root.sort();

    return this;
  }

  initTreeNode(node, data) {
    for (let i = 0, len = data.length; i < len; i++) {
      let _data = data[i];
      let child = new TreeNode(_data);
      node.addChildren(child, true);

      let children = _data[this.option.childrenKey];
      if (children && children.length > 0) {
        this.initTreeNode(child, children);
      }
    }
  }

  update() {
    // let children = this.root.children.concat([]);
    // this.root = new TreeNode();
    // this.root.isRoot = true;
    // this.root.children = children;
    // this.root.sort();
  }

  static traverse(cb, nodes = null) {
    if (!nodes) return;

    let shouldStop = false;

    const result = [];

    for (let nodeInd = 0; nodeInd < nodes.length; nodeInd++) {
      const node = nodes[nodeInd];
      shouldStop = cb(node, nodes) === false;
      result.push(node);

      if (shouldStop) break;

      if (node.children) {
        shouldStop = this.traverse(cb, node.children) === false;
        if (shouldStop) break;
      }
    }

    return !shouldStop ? result : false;
  }

  static findNodeById(nodeModels, id) {
    let node = null;
    TreeTableStore.traverse((nodeModel) => {
      if (nodeModel.node._id === id) {
        node = nodeModel.node;
        return true;
      }
    }, nodeModels);

    return node;
  }

  static getNodeSiblings(nodes, path) {
    if (path.length === 1) return nodes;
    return TreeTableStore.getNodeSiblings(
      nodes[path[0]].children,
      path.slice(1)
    );
  }

  static getNodeByPath(nodes, path) {
    const ind = path.slice(-1)[0];
    let node = TreeTableStore.getNodeSiblings(nodes, path);
    return node ? node[ind] : node;
  }
}
