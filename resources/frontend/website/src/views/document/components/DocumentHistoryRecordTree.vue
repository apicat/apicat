<template>
    <div class="h-full flex flex-col">
        <div class="overflow-x-scroll scroll-content flex-auto" ref="dir">
            <ac-tree
                :data="documentHistoryRecordTree"
                class="bg-transparent"
                node-key="id"
                empty-text=""
                ref="treeIns"
                :expand-on-click-node="false"
                :props="{ children: 'sub_nodes', label: 'title', class: customNodeClass, isLeaf: customNodeLeaf }"
            >
                <template #default="{ node, data }">
                    <div class="el-tree-node__bg"></div>

                    <div class="flex justify-between ac-tree-node">
                        <div class="ac-tree-node__main" @click="handleTreeNodeClick(node, data, $event)">
                            <div class="ac-doc-node" :class="{ 'is-active': data.isCurrent }" :id="'tree_node_' + data.id">
                                <img v-if="data.isLeaf" class="ac-doc-node__icon" :src="createDocIcon" />
                                <span class="ac-doc-node__label" :title="data.title">{{ data.title }}</span>
                            </div>
                        </div>
                    </div>
                </template>
            </ac-tree>
        </div>
    </div>
</template>

<script lang="tsx" setup>
    import { useRoute, useRouter } from 'vue-router'
    import { onMounted, ref, watch } from 'vue'
    import { useDocumentStore } from '@/stores/document'
    import { storeToRefs } from 'pinia'
    import createDocIcon from '@/assets/image/doc-common@2x.png'
    import { traverseTree } from '@natosoft/shared'
    import { DOCUMENT_TYPES } from '@/common/constant'
    import AcTree from './AcTree'
    import { DOCUMENT_HISTORY_DETAIL_NAME } from '@/router/constant'
    import scrollIntoView from 'smooth-scroll-into-view-if-needed'
    import { hideLoading } from '@/hooks/useLoading'

    const $route = useRoute()
    const $router = useRouter()
    const { currentRoute } = $router
    const { params, query } = $route
    const { project_id_public, doc_id } = params

    const documentStore = useDocumentStore()
    const { documentHistoryRecordTree } = storeToRefs(documentStore)

    const treeIns: any = ref(null)
    const dir: any = ref(null)

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
        const { id } = currentRoute.value.params
        // 同一篇文档，且为详情页，不进行任何操作
        if (source.id === parseInt(id as string, 10)) {
            return
        }

        activeNode(source.id)

        $router.push({
            name: DOCUMENT_HISTORY_DETAIL_NAME,
            params: { ...params, id: source.id },
            query,
        })
    }

    const customNodeClass = (data: any) => (data.isLeaf ? 'is-doc' : 'is-dir')

    const customNodeLeaf = (data: any) => data.type === DOCUMENT_TYPES.DOC

    // 文档选中切换
    const activeNode = (rId?: any) => {
        traverseTree(
            (item: any) => {
                if (item.isCurrent) {
                    item.isCurrent = false
                    return false
                }
            },
            documentHistoryRecordTree.value as any,
            { subKey: 'sub_nodes' }
        )

        const { id } = currentRoute.value.params
        const activeNodeId = rId || id

        const _id = parseInt(activeNodeId as string, 10)
        const node = treeIns.value?.getNode(_id)

        if (node) {
            node.data.isCurrent = true
            treeIns.value?.setCurrentKey(_id)
        }
        // scrollIntoView
        const el = document.querySelector('#tree_node_' + _id)
        el && scrollIntoView(el, { scrollMode: 'if-needed' })
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
            documentHistoryRecordTree.value as any,
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
                documentHistoryRecordTree.value,
                { subKey: 'sub_nodes' }
            )

            // 存在文档
            if (node) {
                params.id = node.key
                activeNode(node.key)
            } else {
                hideLoading()
            }

            const dirDom = dir.value as any
            dirDom?.scrollTo(0, 0)

            $router.replace({ name: DOCUMENT_HISTORY_DETAIL_NAME, params, query })
        }
    }

    const getDocTreeList = async () => {
        const tree = await documentStore.getDocumentHistoryRecordList(project_id_public, doc_id)
        if (!tree || !tree.length) {
            hideLoading()
        }
    }

    onMounted(async () => {
        await getDocTreeList()
        params.id ? activeNode() : reactiveNode()
    })
</script>
