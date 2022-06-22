<template>
    <div class="SortableTreeTableRow" :data-path="model.pathStr" :id="node._id">
        <div class="flex-row">
            <div class="td drag justify-content-center"><i class="editor-font editor-drag drag_btn" title="排序"></i></div>
            <div class="td name">
                <div v-for="gap in gaps" :key="gap" class="SortableTreeGap"></div>
                <i v-if="model.children && model.children.length" :class="expandIconClass" @click="onToggleIconClick"></i>

                <slot v-if="!readonly" name="paramName" :row="node">
                    <input v-model="node.name" type="text" placeholder="参数名称" maxlength="100" />
                </slot>

                <span v-else class="copy_text">{{ node.name }}</span>
            </div>

            <div class="td type">
                <select v-model="node.type" :disabled="readonly">
                    <option v-for="type in paramTypes" :value="type.value" :key="type.value">
                        {{ type.text }}
                    </option>
                </select>
            </div>

            <div class="td must justify-content-center">
                <el-checkbox v-if="!readonly" v-model="node.is_must" />
                <span v-if="readonly">{{ node.is_must ? '是' : '否' }}</span>
            </div>

            <div class="td value">
                <input type="text" v-if="!readonly" placeholder="默认值" v-model="node.default_value" :readonly="readonly" maxlength="255" />
                <span v-if="readonly">{{ node.default_value }}</span>
            </div>

            <div class="td desc">
                <input type="text" v-if="!readonly" placeholder="参数说明" v-model="node.description" :readonly="readonly" maxlength="255" />
                <span v-if="readonly">{{ node.description }}</span>
            </div>

            <div class="td mock" @click="onEditMockRuleClick">
                <span :title="node.mock_rule">{{ node.mock_rule }}</span>
            </div>

            <div class="td operations" v-if="!readonly">
                <el-icon :class="{ 'icon-disable': !isAllowAppendSubParams }" title="添加子参数" @click="onAddParamBtnClick"><plus /> </el-icon>
                <el-icon :class="allowAddParamClass" title="添加常用参数" @click="onAddApiParam(node)"><tickets /> </el-icon>
                <el-icon title="删除参数" @click="onRemoveParamBtnClick"><delete /> </el-icon>
            </div>
        </div>

        <div v-if="isShowChild" class="RSortableWrapper" :data-pid="node._id">
            <ResponseParamTreeTableRow
                v-for="(item, index) in model.children"
                :key="item.node._id + index"
                :model="item"
                :editor="editor"
                :readonly="readonly"
                :group="group"
            >
                <template #paramName="{ row }">
                    <slot name="paramName" :row="row"></slot>
                </template>
            </ResponseParamTreeTableRow>
        </div>
    </div>
</template>
<script>
    import { PARAM_TYPES } from '../common/constants'
    import { ElCheckbox, ElIcon } from 'element-plus'
    import { Plus, Tickets, Delete } from '@element-plus/icons-vue'
    import Sortable from 'sortablejs'
    import { $emit } from '@natosoft/shared'

    export default {
        name: 'ResponseParamTreeTableRow',
        inject: ['$treeTable'],
        components: {
            ElCheckbox,
            ElIcon,
            Plus,
            Tickets,
            Delete,
        },
        props: {
            editor: {
                type: Object,
                default() {
                    return null
                },
            },
            group: {
                type: String,
                default: '',
            },
            model: {
                type: Object,
                default: () => ({}),
            },
            readonly: {
                type: Boolean,
                default() {
                    return false
                },
            },
        },
        computed: {
            isShowChild() {
                return this.model.children && this.model.children.length && this.node.expand
            },
            expandIconClass() {
                return ['handle_expand', this.node.expand ? 'el-icon-arrow-down' : 'el-icon-arrow-right']
            },

            gaps() {
                return new Array(Math.max(this.model.level - 1, 0)).fill('gap_').map((v, i) => v + i)
            },

            isRoot() {
                return this.model.level === 1
            },

            node() {
                return this.model.node || {}
            },

            level() {
                return this.model.level
            },

            isAllowAppendSubParams: function () {
                return !((this.node.type !== 4 && this.node.type !== 5 && this.node.type !== 8) || this.level >= 10)
            },

            allowAddParamClass() {
                return [
                    {
                        'icon-disable': this.isAllowAddParam,
                    },
                ]
            },
        },
        watch: {
            'model.node.type': function (newVal, oldVal) {
                if (!this.isAllowAppendSubParams) {
                    this.node.sub_params = []
                }
                this.node.mock_rule = this.$treeTable.getDefaultMockRule(this.node, oldVal)
            },

            'model.node.name': function () {
                this.isAllowAddParam = this.isAllowAddCommonParam(this.node)
            },

            isShowChild: function () {
                this.isShowChild ? this.initSortable() : this.hasDestroyedSortable() ? this.destroyedSortable() : null
            },
        },
        data() {
            return {
                paramTypes: PARAM_TYPES.TYPES.concat([]),
                isAllowAddParam: false,
            }
        },
        methods: {
            onEditMockRuleClick() {
                const root = this.getRoot()
                root.$emit('mock-rule', root, this.model)
            },

            isAllowAddCommonParam(node) {
                if (!node.name.trim()) {
                    return true
                }

                if (!this.editor.commonParamsPopper) {
                    return false
                }

                return this.editor.commonParamsPopper.hasParam(node.name)
            },

            onToggleIconClick() {
                this.node.expand = !this.node.expand
            },

            getRoot() {
                if (this.isRoot) return this
                return this.getParent().getRoot()
            },

            getParent() {
                return this.$parent
            },

            onAddParamBtnClick() {
                if (!this.isAllowAppendSubParams) {
                    return
                }
                this.node.expand = true
                $emit(this.getRoot(), 'add-param', this.model)
            },

            onAddApiParam(node) {
                if (this.isAllowAddCommonParam(node)) {
                    return
                }
                $emit(this.getRoot(), 'add-api-param', this.model)
            },

            onRemoveParamBtnClick() {
                $emit(this.getRoot(), 'remove-param', this.model)
            },

            initSortable() {
                this.$nextTick(() => {
                    this.sortIns = new Sortable(document.querySelector("[data-pid='" + this.node._id + "']"), {
                        group: this.group,
                        animation: 150,
                        handle: '.drag_btn',
                        supportPointer: false,
                        fallbackTolerance: 5,
                        fallbackOnBody: true,
                        swapThreshold: 0.65,
                        forceFallback: this.$treeTable.isSafari || false,
                        fallbackClass: 'sortable-fallback',
                        onEnd: (e) => this.$treeTable.handleSortRows(e),
                    })
                })
            },

            hasDestroyedSortable() {
                return !this.model.children || !this.model.children.length
            },

            destroyedSortable() {
                this.sortIns && this.sortIns.destroy()
                this.sortIns = null
            },
        },

        mounted() {
            this.isShowChild && this.initSortable()
            this.isAllowAddParam = this.isAllowAddCommonParam(this.node)
            this.editor.commonParamsPopper.on(() => {
                this.isAllowAddParam = this.isAllowAddCommonParam(this.node)
            })
        },

        unmounted() {
            this.destroyedSortable()
        },
    }
</script>
