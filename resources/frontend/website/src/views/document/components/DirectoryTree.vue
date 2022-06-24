<template>
    <div class="h-full flex flex-col">
        <div class="ac-doc-catalog">
            <h3 class="text-base font-medium">目录</h3>
            <div>
                <el-icon class="cursor-pointer text-zinc-500" :class="{ 'mr-5': isManager || isDeveloper }" @click="onSearchIconClick"><Search /></el-icon>
                <el-icon v-if="isManager || isDeveloper" class="cursor-pointer text-zinc-500" @click="onRootMoreIconClick"><Plus /></el-icon>
            </div>
        </div>
        <div class="overflow-x-scroll scroll-content flex-auto" ref="dir">
            <ac-tree
                :data="apiDocTree"
                class="bg-transparent"
                node-key="id"
                empty-text=""
                :draggable="isManager || isDeveloper"
                ref="treeIns"
                :expand-on-click-node="false"
                :props="{ children: 'sub_nodes', label: 'title', class: customNodeClass, isLeaf: customNodeLeaf }"
                :allow-drop="allowDrop"
                @node-drag-start="onMoveNodeStart"
                @node-drop="onMoveNode"
            >
                <template #default="{ node, data }">
                    <div class="el-tree-node__bg"></div>

                    <div class="flex justify-between ac-tree-node" :class="{ 'is-editable': data.isEditable }">
                        <div class="ac-tree-node__main" @click="handleTreeNodeClick(node, data, $event)">
                            <div class="ac-doc-node" :class="{ 'is-active': data.isCurrent }">
                                <img v-if="data.isLeaf" class="ac-doc-node__icon" :src="createDocIcon" />
                                <span class="ac-doc-node__label" v-show="!data.isEditable" :title="data.title">{{ data.title }}</span>
                                <input
                                    type="text"
                                    ref="renameInput"
                                    class="ac-doc-node__input el-input el-input__inner"
                                    :id="'tree_input_' + data.id"
                                    v-if="data.isEditable"
                                    v-model="data.title"
                                    @keyup.enter="onEnterKeyUp"
                                    :maxlength="data.isLeaf ? 255 : 50"
                                    @blur="setUnEditable($event, data)"
                                />
                            </div>
                        </div>
                        <div class="ac-tree-node__more" :class="{ active: data.id === activeMoreNodeId }" v-if="isManager || isDeveloper">
                            <el-icon v-show="!data.isLeaf" @click="onMoreIconClick($event, node, data, 'DIR_NEW_TYPE')"><plus /></el-icon>
                            <span class="mx-1" />
                            <el-icon v-show="!data.isLeaf" @click="onMoreIconClick($event, node, data, 'DIR_OPERATE_TYPE')"><more-filled /></el-icon>
                            <el-icon v-show="data.isLeaf" @click="onMoreIconClick($event, node, data, 'DOC_OPERATE_TYPE')"><more-filled /></el-icon>
                        </div>
                    </div>
                </template>
            </ac-tree>
        </div>
    </div>

    <SearchDocumentPopover ref="searchDocumentPopoverRef" />
</template>

<script lang="tsx">
    import { Plus, MoreFilled, Search } from '@element-plus/icons-vue'
    import DirectoryPopper, { NEW_MENUS } from './DirectoryPopper'
    import { useRoute, useRouter } from 'vue-router'
    import { computed, nextTick, onMounted, ref, defineComponent, inject, watch } from 'vue'
    import { useDocumentStore, extendDocTreeFeild } from '@/stores/document'
    import { storeToRefs } from 'pinia'
    import createDocIcon from '@/assets/image/doc-common@2x.png'
    import { memoize, debounce } from 'lodash-es'
    import { traverseTree } from '@natosoft/shared'
    import { DOCUMENT_TYPES } from '@/common/constant'
    import { renameDir, deleteDir, createDir, sortTree } from '@/api/dir'
    import { createDoc, createHttpDoc, renameDoc, deleteDoc, copyDoc, API_DOCUMENT_IMPORT_ACTION_MAPPING } from '@/api/document'
    import NProgress from 'nprogress'
    import { ElMessage as $Message } from 'element-plus'
    import { AsyncMsgBox } from '@/components/AsyncMessageBox'
    import { API_SINGLE_EXPORT_ACTION_MAPPING } from '@/api/exportFile'
    import { hideLoading } from '@/hooks/useLoading'
    import AcTree from './AcTree'
    import SearchDocumentPopover from './SearchDocumentPopover.vue'
    import { useProjectStore } from '@/stores/project'
    import { DOCUMENT_DETAIL_NAME } from '@/router/constant'

    export default defineComponent({
        components: {
            AcTree,
            Plus,
            Search,
            MoreFilled,
            SearchDocumentPopover,
        },

        setup() {
            const documentShareModal: any = inject('documentShareModal')
            const documentImportModal: any = inject('documentImportModal')
            const projectExportModal: any = inject('projectExportModal')

            const $route = useRoute()
            const $router = useRouter()
            const { currentRoute } = $router
            const { params } = $route
            const { project_id } = params
            const documentStore = useDocumentStore()
            const { apiDocTree } = storeToRefs(documentStore)
            const projectStore = useProjectStore()
            const { isManager, isDeveloper } = storeToRefs(projectStore)

            const newMenus: any = NEW_MENUS
            const activeMoreNodeId = ref(null)
            const renameInput: any = ref(null)
            const treeIns: any = ref(null)

            const searchDocumentPopoverRef: any = ref(null)
            const dir: any = ref(null)

            let index = 1
            const oldDraggingNodeInfo: any = null

            // 是否为详情页
            const isDetailPage = computed(() => currentRoute.value.name === DOCUMENT_DETAIL_NAME)

            // 启动切换文档选中
            watch(
                () => $route.params.node_id,
                () => activeNode()
            )

            const handleTreeNodeClick = (node: any, source: any, e: any) => {
                if (e.target.tagName === 'INPUT') {
                    return
                }

                // 文档点击
                if (source.isLeaf) {
                    onDocumentClick(source)
                    return
                }
                // 目录点击
                node.expanded = !node.expanded
            }

            const onDocumentClick = (source: any) => {
                const { node_id } = currentRoute.value.params
                // 同一篇文档，且为详情页，不进行任何操作
                if (isDetailPage.value && source.id === parseInt(node_id as string, 10)) {
                    return
                }
                activeNode(source.id)
                $router.push({ name: DOCUMENT_DETAIL_NAME, params: { ...params, node_id: source.id } })
            }

            const onMoreIconClick = (e: any, node: any, source: any, type: any) => {
                e.stopPropagation()
                activeMoreNodeId.value = source.id
                DirectoryPopper.show(type, e.target, { node, source, isDetailPage: isDetailPage.value }, dir.value)
            }

            const onRootMoreIconClick = (e: any) => {
                onMoreIconClick(e, treeIns.value?.root, apiDocTree.value, 'DIR_NEW_TYPE')
            }

            const getTreeMaxDepth = memoize(function (node) {
                let maxLevel = 0
                traverseTree(
                    (item: any) => {
                        if (!item.isLeaf) {
                            maxLevel++
                        }
                    },
                    [node] as any,
                    { subKey: 'sub_nodes' }
                )

                return maxLevel
            })

            const allowDrop = (draggingNode: any, dropNode: any, type: any) => {
                const { data: dropNodeData } = dropNode
                const { data: draggingNodeData } = draggingNode

                // 不允许拖放在文件中
                if (dropNodeData.isLeaf && type === 'inner') {
                    return false
                }

                // 拖动目录时
                if (!draggingNodeData.isLeaf && !dropNodeData.isLeaf) {
                    return getTreeMaxDepth(draggingNodeData) + dropNode.level <= 5
                }

                return true
            }

            const customNodeClass = (data: any) => (data.isLeaf ? 'is-doc' : 'is-dir')

            const customNodeLeaf = (data: any) => data.type === DOCUMENT_TYPES.DOC
            // 重命名功能
            const inputFocus = async () => {
                await nextTick()
                renameInput.value?.focus()
                renameInput.value?.setSelectionRange(0, renameInput.value?.value.length)
            }

            const onEnterKeyUp = (e: any) => {
                e.target && e.target.blur()
            }

            const setUnEditable = debounce(function (e, source) {
                source.isEditable = false

                // 进行数据还原
                if (!e.target.value && source._oldName) {
                    source.title = source._oldName
                    return
                }

                const isDir = source.type === DOCUMENT_TYPES.DIR
                const action = isDir ? renameDir : renameDoc
                const param = { project_id, title: source.title } as any
                param[isDir ? 'node_id' : 'doc_id'] = source.id

                action(param)
                    .then(() => {
                        delete source._oldName
                    })
                    .catch(() => {
                        if (source._oldName) source.title = source._oldName
                    })
            }, 200)

            // 文档选中切换
            const activeNode = (nodeId?: any) => {
                traverseTree(
                    (item: any) => {
                        if (item.isCurrent) {
                            item.isCurrent = false
                            return false
                        }
                    },
                    apiDocTree.value as any,
                    { subKey: 'sub_nodes' }
                )

                const { node_id } = currentRoute.value.params
                const activeNodeId = nodeId || node_id

                const id = parseInt(activeNodeId as string, 10)
                const node = treeIns.value?.getNode(id)

                if (node) {
                    node.data.isCurrent = true
                    treeIns.value?.setCurrentKey(id)
                }
            }

            const reactiveNode = () => {
                if (!treeIns.value) {
                    return
                }

                let hasCurrent = false
                traverseTree(
                    (item: any) => {
                        if (item.isCurrent) {
                            hasCurrent = true
                            return false
                        }
                    },
                    apiDocTree.value as any,
                    { subKey: 'sub_nodes' }
                )

                // 没有选中文档时，进行自动切换
                if (!hasCurrent) {
                    let node: any = null
                    traverseTree(
                        (item: any) => {
                            let _node = treeIns.value.getNode(item.id)
                            if (_node && _node.data.isLeaf) {
                                node = _node
                                return false
                            }
                        },
                        apiDocTree.value,
                        { subKey: 'sub_nodes' }
                    )

                    let params = { project_id } as any

                    // 存在文档
                    if (node) {
                        params.node_id = node.key
                        activeNode(node.key)
                    }

                    const dirDom = dir.value as any
                    dirDom?.scrollTo(0, 0)
                    $router.replace({ name: DOCUMENT_DETAIL_NAME, params })
                }
            }

            const getDocTreeList = async () => {
                const tree = await documentStore.getApiDocTree(project_id as string)
                hideLoading()
            }

            const onSearchIconClick = (e: any) => {
                searchDocumentPopoverRef.value?.show(e.currentTarget)
            }

            onMounted(async () => {
                await getDocTreeList()
                currentRoute.value.params.node_id ? activeNode() : reactiveNode()
            })

            return {
                isManager,
                isDeveloper,
                router: $router,
                index,
                oldDraggingNodeInfo,

                project_id,

                apiDocTree,
                renameInput,
                dir,
                activeMoreNodeId,
                treeIns,
                newMenus,
                documentShareModal,

                createDocIcon,

                handleTreeNodeClick,
                onRootMoreIconClick,
                onMoreIconClick,
                onSearchIconClick,
                allowDrop,
                customNodeClass,
                customNodeLeaf,

                onEnterKeyUp,
                setUnEditable,
                inputFocus,

                activeNode,
                reactiveNode,

                documentImportModal,
                projectExportModal,
                searchDocumentPopoverRef,
            }
        },

        mounted() {
            DirectoryPopper.onPopperItemClick = ({ menu, node, source }) => this.onPopperItemClick(menu, node, source)
            DirectoryPopper.onPopperHide = () => {
                this.activeMoreNodeId = null
            }
        },

        methods: {
            onPopperItemClick(menu: any, node?: any, data?: any) {
                const fn: any = this[menu.onClick]
                fn && fn(node || this.treeIns.root, data || { sub_nodes: this.apiDocTree, id: 0 }, menu)
            },

            onRenameBtnClick(node: any, source: any) {
                source._oldName = source.title
                source.isEditable = true
                this.inputFocus()
            },

            onDeleteBtnClick(node: any, source: any) {
                let isDir = !source.isLeaf
                AsyncMsgBox({
                    title: '删除提示',
                    content: (
                        <div class="break-all">
                            确定删除「{source.title}」{isDir ? '分类' : '文件'}吗？
                        </div>
                    ),
                    onOk: () => {
                        const isDir = source.type === DOCUMENT_TYPES.DIR
                        const action = isDir ? deleteDir : deleteDoc
                        let param = { project_id: this.project_id } as any
                        param[isDir ? 'node_id' : 'doc_id'] = source.id

                        return action(param).then((res: any) => {
                            $Message.success(res.msg || '删除成功！')
                            this.treeIns.remove(node)
                            this.reactiveNode()
                        })
                    },
                })
            },

            onCreateDirBtnClick(node: any, source: any) {
                let data = {
                    project_id: this.project_id,
                    name: '新建分类' + this.index++,
                } as any

                if (source.id) {
                    data.parent_id = source.id
                    data.pid = source.id
                }

                NProgress.start()
                createDir(data)
                    .then((res) => {
                        const data = extendDocTreeFeild(res.data, DOCUMENT_TYPES.DIR)

                        // root
                        if (!node.parent && node.level === 0) {
                            source.unshift(data)
                        } else {
                            if (!source.sub_nodes || !source.sub_nodes.length) {
                                this.treeIns.append(data, node)
                            } else {
                                this.treeIns.insertBefore(data, source.sub_nodes[0])
                            }
                        }

                        this.$nextTick(() => {
                            const parentNode = this.treeIns.getNode(source)
                            parentNode && (parentNode.expanded = true)
                            const current = this.treeIns.getNode(res.data)
                            current.data.isEditable = true
                            setTimeout(() => this.inputFocus(), 100)
                        })
                    })
                    .finally(() => {
                        NProgress.done()
                    })
            },

            createDoc(node: any, source: any) {
                let param = { project_id: this.project_id, parent_id: node.key, pid: node.key }
                NProgress.start()
                createDoc(param)
                    .then(({ data }) => {
                        this.treeIns.append(extendDocTreeFeild(data), node)
                        this.$nextTick(() => {
                            this.treeIns.setCurrentKey(data.id)
                            const parentNode = this.treeIns.getNode(source)
                            parentNode && (parentNode.expanded = true)

                            this.router.push({
                                name: 'document.api.edit',
                                params: { project_id: this.project_id, node_id: data.id },
                                query: { isNew: true } as any,
                            })
                            this.activeNode(data.id)
                        })
                    })
                    .finally(() => {
                        NProgress.done()
                    })
            },

            onMoveNodeStart(draggingNode: any) {
                const oldParent = draggingNode.parent

                this.oldDraggingNodeInfo = {
                    oldPid: oldParent.id === 0 ? null : oldParent.key,
                    oldChildIds: oldParent.childNodes.filter((item: any) => item.id !== draggingNode.id).map((item: any) => item.key),
                }
            },

            onMoveNode(draggingNode: any, dropNode: any, dropType: string) {
                if (!this.oldDraggingNodeInfo) {
                    return
                }

                const { oldPid, oldChildIds } = this.oldDraggingNodeInfo

                let isSeamLevel = oldPid === dropNode.parent.id && dropType !== 'inner'
                const newParent = this.treeIns.getNode(draggingNode.data).parent

                const newPid = newParent.id === 0 ? null : newParent.key
                const newChildIds = newParent.childNodes.map((item: any) => item.key)

                const sortData = {
                    oldPid: isSeamLevel ? newPid : oldPid,
                    oldChildIds: isSeamLevel ? newChildIds : oldChildIds,
                    newPid,
                    newChildIds,
                }

                this.oldDraggingNodeInfo = null

                sortTree({
                    new_pid: sortData.newPid || 0,
                    new_node_ids: sortData.newChildIds,
                    old_pid: sortData.oldPid || 0,
                    old_node_ids: sortData.oldChildIds,
                    project_id: this.project_id,
                })
            },

            onCopyBtnClick(node: any, source: any) {
                NProgress.start()
                copyDoc({
                    project_id: this.project_id,
                    doc_id: source.id,
                })
                    .then(({ data }) => {
                        this.treeIns.insertAfter(extendDocTreeFeild(data), node)
                        this.router.push({ name: 'document.api.edit', params: { project_id: this.project_id, node_id: data.id } })
                    })
                    .finally(() => {
                        NProgress.done()
                    })
            },

            onShareBtnClick(node: any, source: any) {
                this.documentShareModal.show({
                    docId: source.id,
                    nodeId: source.id,
                })
            },

            onImportBtnClick(node: any, source: any) {
                this.documentImportModal.show(
                    {
                        project_id: this.project_id,
                        parent_id: source.id ? source.id : 0,
                    },
                    API_DOCUMENT_IMPORT_ACTION_MAPPING
                )
            },

            onExportBtnClick(node: any, source: any) {
                this.projectExportModal.title = '导出文档'
                this.projectExportModal.show({ project_id: this.project_id, doc_id: source.id }, API_SINGLE_EXPORT_ACTION_MAPPING)
            },

            // 更新文档标题
            updateTreeNode(id: any, newNode = { title: '' }) {
                let node = this.treeIns.getNode(id)
                if (node && node.data.title) {
                    node.data.title = newNode.title || node.title
                }
            },

            createHttpDoc(node: any, source: any) {
                let param = { project_id: this.project_id, parent_id: node.key, pid: node.key }
                NProgress.start()
                createHttpDoc(param)
                    .then(({ data }) => {
                        this.treeIns.append(extendDocTreeFeild(data), node)
                        this.$nextTick(() => {
                            const parentNode = this.treeIns.getNode(source)
                            parentNode && (parentNode.expanded = true)

                            this.router.push({
                                name: 'document.api.edit',
                                params: { project_id: this.project_id, node_id: data.id },
                                query: { isNew: true } as any,
                            })
                        })
                    })
                    .finally(() => {
                        NProgress.done()
                    })
            },
        },
    })
</script>
