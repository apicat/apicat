<template>
    <div class="SortableTreeTable-scroll">
        <div class="SortableTreeTable" :id="id">
            <div class="SortableTreeTableHeader">
                <div class="flex-row">
                    <div class="th drag" v-html="'&nbsp;'"></div>
                    <div class="th name">参数名称</div>
                    <div class="th type">参数类型</div>
                    <div class="th must ivu-row-flex-center">必传</div>
                    <div class="th value">默认值</div>
                    <div class="th desc">参数说明</div>
                    <div v-if="!readonly" class="th operations" v-html="'&nbsp;'"></div>
                </div>
            </div>

            <div class="SortableTreeTableBody RSortableWrapper" ref="table_body_wrapper">
                <TreeTableRow
                    v-for="node in nodes"
                    :key="node.node._id"
                    :group="id"
                    :model="node"
                    :editor="editor"
                    :readonly="readonly"
                    @add-param="onAddSubParam"
                    @remove-param="onRemoveParam"
                    @add-api-param="onAddCommonParam"
                >
                    <template #paramName="{ row }">
                        <el-autocomplete
                            ref="autocomplete"
                            :fetch-suggestions="querySuggestionsList"
                            @select="(val) => onParamItemClick(row, val)"
                            v-model="row.name"
                            placeholder="参数名称"
                            :maxlength="100"
                        >
                            <template #default="{ item }">
                                <div class="ac-complete-item">
                                    <div class="ac-complete-item-content">
                                        {{ item.value }}
                                    </div>
                                    <el-icon @click.stop="onDeleteParamBtnClick($event, item.value)"><delete></delete></el-icon>
                                </div>
                            </template>
                        </el-autocomplete>
                    </template>
                </TreeTableRow>
            </div>
            <button v-if="!readonly" class="add-root-param" @click="onAddRootParamBtnClick">添加</button>
        </div>
    </div>
</template>

<script>
    import { ElAutocomplete, ElIcon } from 'element-plus'
    import { Delete } from '@element-plus/icons-vue'
    import shortid from 'shortid'
    import Sortable from 'sortablejs'
    import { removeNode, insertNodeAt } from './utils'
    import TreeTableRow from './ParamTreeTableRow.vue'
    import TreeTableStore from './TreeTableStore'
    import { $emit } from '@ac/shared'

    export default {
        name: 'ParamTreeTable',
        components: {
            TreeTableRow,
            ElAutocomplete,
            ElIcon,
            Delete,
        },
        props: {
            editor: {
                type: Object,
                default() {
                    return null
                },
            },
            readonly: {
                type: Boolean,
                default() {
                    return false
                },
            },
            data: {
                type: Array,
                default() {
                    return []
                },
            },
            expand: {
                type: Boolean,
                default() {
                    return true
                },
            },
            paramsTip: {
                type: Array,
                default() {
                    return []
                },
            },
        },

        data() {
            return {
                id: 'ParamTreeTable_' + shortid(),
                isSafari:
                    navigator.vendor &&
                    navigator.vendor.indexOf('Apple') > -1 &&
                    navigator.userAgent &&
                    navigator.userAgent.indexOf('CriOS') === -1 &&
                    navigator.userAgent.indexOf('FxiOS') === -1,
                nodes: this.createNodes(this.data),
                commonParamsManager: null,
            }
        },

        provide() {
            return {
                $treeTable: this,
            }
        },

        watch: {
            data: {
                deep: true,
                handler: function (newVal) {
                    this.nodes = this.createNodes(newVal)
                },
            },
        },

        methods: {
            createNodes(data) {
                return this.getNodes(data)
            },

            getNodes(nodeModels, parentPath = [], parent) {
                return nodeModels.map((nodeModel, ind) => {
                    const nodePath = parentPath.concat(ind)

                    nodeModel._id = nodeModel._id || shortid()

                    nodeModel['expand'] = nodeModel.expand !== undefined ? nodeModel.expand : this.expand
                    // this.$set(nodeModel, 'expand', nodeModel.expand !== undefined ? nodeModel.expand : this.expand)

                    return this.getNode(nodePath, nodeModel, parent)
                })
            },

            getNode(path, nodeModel = null, parent = null) {
                if (!nodeModel) return null
                return {
                    parent,
                    node: nodeModel,
                    children: nodeModel.sub_params ? this.getNodes(nodeModel.sub_params, path, nodeModel) : [],
                    path: path,
                    pathStr: JSON.stringify(path),
                    level: path.length,
                }
            },

            getNodeSiblings(nodes, path) {
                if (path.length === 1) return nodes
                return this.getNodeSiblings(nodes[path[0]].children, path.slice(1))
            },

            getNodeByPath(path) {
                const ind = path.slice(-1)[0]
                let node = this.getNodeSiblings(this.nodes, path)
                return node ? node[ind] : node
            },

            onAddRootParamBtnClick() {
                this.data.push(this.generateSubParam())
            },

            onAddSubParam(model) {
                const { node } = model
                if (!node.sub_params) {
                    node.sub_params = []
                }
                node.sub_params.push(this.generateSubParam())

                $emit(this, 'add-param', node)
            },

            onRemoveParam(model) {
                const { node, path } = model

                let len = path.length - 1
                let index = path[len]
                const parent = model.parent ? model.parent.sub_params : this.data
                parent && parent.splice(index, 1)

                $emit(this, 'remove-param', node)
            },

            generateSubParam() {
                return {
                    name: '',
                    type: 1,
                    is_must: false,
                    default_value: '',
                    description: '',
                    sub_params: [],
                }
            },

            handleSortRows(event) {
                const { item, from, to, oldIndex, newIndex } = event
                // console.log(item, from, to, oldIndex, newIndex);
                // console.log(from === to ? "同级" : "跨级");

                // 同级
                if (from === to) {
                    let pathStr = item.getAttribute('data-path')

                    pathStr = JSON.parse(pathStr)

                    if (oldIndex === newIndex || !pathStr) {
                        return
                    }

                    let node = this.getNodeByPath(pathStr)

                    let arr = this.data

                    if (node.parent) {
                        arr = node.parent.sub_params
                    }
                    let old = arr.splice(oldIndex, 1)
                    old && old.length && arr.splice(newIndex, 0, old[0])
                }
                // 跨级
                else {
                    let fromPid = from.getAttribute('data-pid')
                    let toPid = to.getAttribute('data-pid')

                    let fromNode = TreeTableStore.findNodeById(this.nodes, fromPid)
                    let toNode = TreeTableStore.findNodeById(this.nodes, toPid)
                    let fromNodeArray = fromNode ? fromNode.sub_params : this.data
                    let toNodeArray = toNode ? toNode.sub_params : this.data

                    // console.log("fromNode:", fromNode, "fromNodeArray:", fromNodeArray);
                    // console.log("toNode:", toNode, "toNodeArray:", toNodeArray);

                    // window.fromNodeArray = fromNodeArray;
                    // window.toNodeArray = toNodeArray;
                    // window.oldIndex = oldIndex;
                    // window.newIndex = newIndex;

                    removeNode(event.item)
                    insertNodeAt(event.from, event.item, event.oldIndex)

                    let nodes = fromNodeArray.splice(oldIndex, 1)
                    nodes && nodes.length && toNodeArray.splice(newIndex, 0, nodes[0])
                }
            },

            initSortable() {
                this.rootSortIns = new Sortable(this.$refs.table_body_wrapper, {
                    group: this.id,
                    animation: 150,
                    handle: '.drag_btn',
                    supportPointer: false,
                    fallbackTolerance: 5,
                    fallbackOnBody: true,
                    swapThreshold: 0.65,
                    forceFallback: this.isSafari,
                    fallbackClass: 'sortable-fallback',
                    onEnd: (e) => this.handleSortRows(e),
                })
            },

            querySuggestionsList(queryString, cb) {
                if (!this.editor || !this.editor.commonParamsManager) {
                    return cb([])
                }
                cb(this.editor.commonParamsManager.queryParams(queryString))
            },

            onParamItemClick(row, { value }) {
                let newNode = null
                if (this.editor.commonParamsManager && (newNode = this.editor.commonParamsManager.getParamByKey(value))) {
                    row.name = newNode.name
                    row.type = newNode.type
                    row.default_value = newNode.default_value
                    row.is_must = newNode.is_must
                    row.description = newNode.description
                }
            },

            onDeleteParamBtnClick(e, key) {
                e.stopPropagation()
                e.preventDefault()
                this.editor.commonParamsManager && this.editor.commonParamsManager.deleteParam(key)

                let target = e.target
                let result = null
                while (target) {
                    if (target.tagName !== 'UL') {
                        target = target.parentElement
                    } else {
                        result = target
                        target = null
                    }
                }

                if (result) {
                    const input = document.querySelector(`[aria-owns="${result.getAttribute('id')}"] input`)

                    input && input.focus()
                }
            },

            onAddCommonParam({ node }) {
                this.editor.commonParamsManager && this.editor.commonParamsManager.addParam(node)
            },
        },

        mounted() {
            this.initSortable()
        },

        unmounted() {
            this.rootSortIns && this.rootSortIns.destroy()
            this.rootSortIns = null
        },
    }
</script>
