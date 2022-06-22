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
                    <div class="th mock">Mock</div>
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
                    @mock-rule="onEditMockRuleClick"
                    @remove-param="onRemoveParam"
                    @add-api-param="onAddCommonParam"
                >
                    <template #paramName="{ row }">
                        <el-input
                            v-model="row.name"
                            placeholder="参数名称"
                            :maxlength="100"
                            @input="onParamNameChange"
                            @focus="onParamNameInputFocus($event, row)"
                        />
                    </template>
                </TreeTableRow>
            </div>
            <button v-if="!readonly" class="add-root-param" @click="onAddRootParamBtnClick">添加</button>
        </div>
    </div>
</template>

<script>
    import shortid from 'shortid'
    import Sortable from 'sortablejs'
    import { isPlainObject } from 'lodash-es'

    import TreeTableRow from './ResponseParamTreeTableRow.vue'
    import TreeTableStore from './TreeTableStore'
    import { ElInput, ElIcon } from 'element-plus'
    import { Delete } from '@element-plus/icons-vue'
    import { insertNodeAt, removeNode, generateArray } from './utils'
    import { $emit } from '@natosoft/shared'

    import { getMockRules, PARAM_TYPES } from '../common/constants'

    const MOCK_RULES = getMockRules()

    export default {
        name: 'ResponseParamTreeTable',
        components: {
            TreeTableRow,
            ElInput,
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
            onParamNameChange(val) {
                this.editor.commonParamsPopper && this.editor.commonParamsPopper.queryParams(val)
            },

            onParamNameInputFocus($event, node) {
                this.editor.commonParamsPopper &&
                    this.editor.commonParamsPopper.show({
                        inputDom: $event.target,
                        node,
                        onParamItemClick: (node, paramKey) => this.onParamItemClick(node, { value: paramKey }),
                    })
            },

            createNodes(data) {
                return this.getNodes(data)
            },

            getNodes(nodeModels, parentPath = [], parent) {
                return nodeModels.map((nodeModel, ind) => {
                    const nodePath = parentPath.concat(ind)

                    nodeModel._id = nodeModel._id || shortid()
                    nodeModel['expand'] = nodeModel.expand !== undefined ? nodeModel.expand : this.expand
                    nodeModel['mock_rule'] = nodeModel.mock_rule !== undefined ? nodeModel.mock_rule : ''

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

            onEditMockRuleClick(vm, model) {
                $emit(this, 'mock-rule', vm, model)
            },

            onAddRootParamBtnClick() {
                // eslint-disable-next-line vue/no-mutating-props
                this.data.push(this.generateSubParam())
            },

            onAddSubParam(model) {
                const { node } = model
                if (!node.sub_params) {
                    node.sub_params = []
                }
                node.sub_params.push(this.generateSubParam(model))

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

            generateSubParam(parentModel) {
                let name = ''

                if (parentModel && parentModel.node.type === PARAM_TYPES.VALUES.ARRAY && parentModel.node.name) {
                    name = parentModel.node.name + generateArray(1)
                }

                const node = {
                    name,
                    type: 1,
                    is_must: false,
                    mock_rule: '',
                    default_value: '',
                    description: '',
                    sub_params: [],
                }

                node.mock_rule = this.getDefaultMockRule(node)

                return node
            },

            handleSortRows(event) {
                const { item, from, to, oldIndex, newIndex } = event
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
                if (!this.editor || !this.editor.commonParamsPopper) {
                    return cb([])
                }
                cb(this.editor.commonParamsPopper.queryParams(queryString))
            },

            onParamItemClick(row, { value }) {
                let newNode = null
                if (this.editor.commonParamsPopper && (newNode = this.editor.commonParamsPopper.getParamByKey(value))) {
                    row.name = newNode.name
                    row.type = newNode.type
                    row.default_value = newNode.default_value
                    row.is_must = newNode.is_must
                    row.description = newNode.description
                }
            },

            onAddCommonParam({ node }) {
                this.editor.commonParamsPopper && this.editor.commonParamsPopper.addParam(node)
            },

            getDefaultMockRule(node, oldType) {
                if (!isPlainObject(node)) {
                    throw new Error('响应参数类型有误！')
                }

                const paramType = PARAM_TYPES.valueOf(node.type).toLowerCase()
                const mockInfo = MOCK_RULES[paramType]

                let defaultRule = ''

                if (mockInfo) {
                    const defaultMockRule = mockInfo.rules

                    // 空参数名称，默认规则
                    if (!node.name) {
                        return this.getMockRuleDefaultValue(defaultMockRule, defaultMockRule[0].name)
                    }

                    // 精准匹配规则
                    if (mockInfo.ruleKeys.indexOf(node.name) !== -1) {
                        return this.getMockRuleDefaultValue(defaultMockRule, node.name)
                    }

                    // 类型检索
                    const len = defaultMockRule.length
                    for (let i = 0; i < len; i++) {
                        const rule = defaultMockRule[i]
                        if (rule.searchKey.indexOf(paramType) !== -1) {
                            defaultRule = rule.default || rule.name
                            break
                        }
                    }

                    // searchkeys 再次精准匹配
                    for (let i = 0; i < len; i++) {
                        const rule = defaultMockRule[i]
                        if (rule.searchKeys.indexOf(node.name) !== -1) {
                            defaultRule = rule.default || rule.name
                            break
                        }
                    }
                }

                // object -> array_object
                if (defaultRule === 'array' && oldType === 5 && node.sub_params && node.sub_params.length >= 2) {
                    defaultRule = 'array_object'
                }

                return defaultRule
            },

            getMockRuleDefaultValue(rules, ruleName) {
                const rule = (rules || []).find((rule) => rule.name === ruleName)

                return rule ? rule.default || rule.name : ''
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
