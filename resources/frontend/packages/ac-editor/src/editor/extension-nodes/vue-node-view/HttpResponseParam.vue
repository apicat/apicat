<template>
    <node-view-wrapper :class="className">
        <el-tabs v-model="activeName">
            <el-tab-pane v-for="item in tabs" :key="item.key" :label="item.key">
                <template #label>
                    <div class="tab-header--editable">
                        <p v-show="!item.isEdit">
                            {{ item.title }}<el-icon @click="onEditTitleIconClick(item)"><edit-pen></edit-pen></el-icon>
                        </p>
                        <input
                            :id="item.id"
                            v-show="item.isEdit"
                            type="text"
                            :value="item.title"
                            @keydown.enter="onEnter"
                            @keydown.esc="onCancel($event, item)"
                            @blur="onChangeTabName($event, item)"
                        />
                    </div>
                </template>
                <ResponseParamTreeTable @mock-rule="onEditMockRuleClick" :editor="editor" :readonly="isReadOnly" :data="item.data" />
            </el-tab-pane>
        </el-tabs>
    </node-view-wrapper>
</template>

<script>
    import { ElTabs, ElTabPane, ElIcon } from 'element-plus'
    import { EditPen } from '@element-plus/icons-vue'
    import ResponseParamTreeTable from '../../../components/ResponseParamTreeTable.vue'
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'
    import { ref, watch, nextTick, toRaw } from 'vue'
    import shortid from 'shortid'

    export default {
        name: 'HttpResponseParam',
        props: ['editor', 'node', 'view', 'getPos', 'isSelected', 'isReadOnly', 'updateAttributes'],
        components: {
            ElTabs,
            ElTabPane,
            ElIcon,
            EditPen,
            ResponseParamTreeTable,
            NodeViewWrapper,
        },
        setup(props) {
            const headerParams = ref(props.node.attrs.response_header.params)
            const bodyParams = ref(props.node.attrs.response_body.params)

            const tabs = ref([
                { id: 'input_' + shortid(), key: 'response_header', data: headerParams, title: props.node.attrs.response_header.title, isEdit: false },
                { id: 'input_' + shortid(), key: 'response_body', data: bodyParams, title: props.node.attrs.response_body.title, isEdit: false },
            ])

            const onEditTitleIconClick = async (item) => {
                item.isEdit = !item.isEdit
                await nextTick()
                const input = document.querySelector('#' + item.id)
                input && input.focus()
            }

            const onEnter = (e) => {
                e.target.blur && e.target.blur()
            }

            const onChangeTabName = (e, item) => {
                item.isEdit = false
                if (!e.target.value.trim()) {
                    return
                }
                item.title = e.target.value
                updateTabTitle(item.key, e.target.value)
            }

            const updateTabTitle = (paramKey, newVal) => {
                let oldData = props.node.attrs[paramKey]
                oldData && props.updateAttributes({ [paramKey]: { ...oldData, title: newVal } })
            }

            watch(
                [headerParams, bodyParams],
                ([newHeaderParams, newBodyParams]) => {
                    const hTitle = props.node.attrs.response_header.title
                    const bTitle = props.node.attrs.response_body.title

                    props.updateAttributes({
                        response_header: { params: toRaw(newHeaderParams), title: hTitle },
                        response_body: { params: toRaw(newBodyParams), title: bTitle },
                    })
                },
                { deep: true }
            )

            return {
                tabs,

                onEnter,
                onChangeTabName,
                onEditTitleIconClick,
            }
        },

        data() {
            return {
                activeName: '0',
            }
        },
        computed: {
            className() {
                return [
                    'http-response-param',
                    {
                        'ProseMirror-selectednode': this.isSelected && !this.isReadOnly,
                        'http-param-readonly': this.isReadOnly,
                    },
                ]
            },
        },
        methods: {
            initActiveTab() {
                const { response_header, response_body } = this.node.attrs

                if (!response_header.params.length && response_body.params.length) {
                    this.activeName = '1'
                }
            },

            onEditMockRuleClick(vm, model) {
                this.editor.mockModel && this.editor.mockModel.show(vm, model)
            },
        },
        created() {
            this.initActiveTab()
        },
    }
</script>
