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
                <ParamTreeTable :editor="editor" :readonly="isReadOnly" :data="item.data" />
            </el-tab-pane>
        </el-tabs>
    </node-view-wrapper>
</template>

<script>
    import { ElTabs, ElTabPane, ElIcon } from 'element-plus'
    import { EditPen } from '@element-plus/icons-vue'
    import ParamTreeTable from '../../../components/ParamTreeTable.vue'
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'
    import { ref, watch, nextTick, toRaw } from 'vue'
    import shortid from 'shortid'

    export default {
        name: 'HttpRequestParam',
        props: ['editor', 'node', 'view', 'getPos', 'isSelected', 'isReadOnly', 'updateAttributes'],
        components: {
            ElTabs,
            ElTabPane,
            ElIcon,
            EditPen,
            ParamTreeTable,
            NodeViewWrapper,
        },
        setup(props) {
            const headerParams = ref(props.node.attrs.request_header.params)
            const bodyParams = ref(props.node.attrs.request_body.params)
            const queryParams = ref(props.node.attrs.request_query.params)

            const tabs = ref([
                { id: 'input_' + shortid(), key: 'request_header', data: headerParams, title: props.node.attrs.request_header.title, isEdit: false },
                { id: 'input_' + shortid(), key: 'request_body', data: bodyParams, title: props.node.attrs.request_body.title, isEdit: false },
                { id: 'input_' + shortid(), key: 'request_query', data: queryParams, title: props.node.attrs.request_query.title, isEdit: false },
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
                [headerParams, bodyParams, queryParams],
                ([newHeaderParams, newBodyParams, newQueryParams]) => {
                    const hTitle = props.node.attrs.request_header.title
                    const bTitle = props.node.attrs.request_body.title
                    const qTitle = props.node.attrs.request_query.title
                    props.updateAttributes({
                        request_header: { params: toRaw(newHeaderParams), title: hTitle },
                        request_body: { params: toRaw(newBodyParams), title: bTitle },
                        request_query: { params: toRaw(newQueryParams), title: qTitle },
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
                    'http-request-param',
                    {
                        'ProseMirror-selectednode': this.isSelected && !this.isReadOnly,
                        'http-param-readonly': this.isReadOnly,
                    },
                ]
            },
            attrs() {
                return this.node.attrs
            },
        },
        methods: {
            initActiveTab() {
                const { request_header, request_body, request_query } = this.node.attrs

                if (!request_header.params.length && request_query.params.length) {
                    this.activeName = '2'
                }

                if (!request_header.params.length && request_body.params.length) {
                    this.activeName = '1'
                }
            },
        },
        created() {
            this.initActiveTab()
        },
    }
</script>
