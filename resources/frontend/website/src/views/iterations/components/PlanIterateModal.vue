<template>
    <el-dialog v-model="isShow" :width="800" :close-on-click-modal="false" title="规划迭代" append-to-body class="show-footer-line">
        <p class="title border-l-[3px] p-2 mb-3 italic">规划本次迭代所涉及的API</p>
        <AcTransferTree
            v-loading="isLoadingForTree"
            height="500px"
            ref="transferTreeRef"
            :defaultProps="defaultProps"
            :from_data="fromData"
            :to_data="toData"
            pid="parent_id"
            filter
            :title="['所有API', '已规划API']"
            @addBtn="onChangeTree"
            @removeBtn="onChangeTree"
        />
        <template v-slot:footer>
            <el-button @click="onCloseBtnClick()"> 取消 </el-button>
            <el-button :loading="isLoading" type="primary" @click="onOkBtnClick"> 确 定 </el-button>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
    import { ref, unref, watch } from 'vue'
    import { planApisToIterate } from '@/api/iterate'
    import { treeList } from '@/api/dir'
    import useApi from '@/hooks/useApi'
    import { flattenDeep, arrayToTree, traverseForMarkSort } from '@/common/helper'
    import { DOCUMENT_TYPES } from '@/common/constant'

    const defaultProps = {
        label: 'title',
        children: 'sub_nodes',
    }

    const emit = defineEmits(['ok'])

    const [isLoadingForTree, execute] = useApi(treeList)
    const [isLoading, executePlanApisToIterate] = useApi(planApisToIterate, { isCatch: false })
    const isShow = ref(false)

    const iterateId = ref(null)

    const transferTreeRef = ref()
    const fromData: any = ref([])
    const toData: any = ref([])
    const expandedKeys: any = ref([])

    const show = (iterate: any) => {
        isShow.value = true
        iterateId.value = iterate.id
    }

    const hide = () => {
        isShow.value = false
    }

    const onCloseBtnClick = () => {
        hide()
    }

    const onOkBtnClick = async () => {
        const iteration_id = unref(iterateId)
        if (iteration_id === null) {
            return
        }
        const node_ids = transferTreeRef.value.getFlattenValues().map((item: any) => item.id)

        // if (!node_ids.length) {
        //     $Message.error('请选择需要规划的API')
        //     return
        // }
        await executePlanApisToIterate({ iteration_id, node_ids })
        emit('ok')
        hide()
    }

    watch(isShow, () => {
        if (!isShow.value) {
            fromData.value = []
            toData.value = []
            expandedKeys.value = []
            iterateId.value = null
            transferTreeRef.value?.clearExpand()
        }
    })

    watch(iterateId, async () => {
        const iteration_id = unref(iterateId)
        if (iteration_id === null) {
            return
        }

        const { data: allData } = await execute({ iteration_id })

        traverseForMarkSort({ sub_nodes: allData })

        const arr = flattenDeep(allData || [], defaultProps.children)

        const toDataArr: any = []
        const removeIndexSet = new Set()

        arr.forEach((item: any) => {
            if (item.selected) {
                item.sub_nodes = []
                toDataArr.push(item)
                // 移除所有文档类型
                if (item.type === DOCUMENT_TYPES.DOC) {
                    removeIndexSet.add(item.id)
                }
            }
        })

        let fromDataArr = arr.filter((item: any) => !removeIndexSet.has(item.id))

        // 清空所有子集，重新组成新树
        fromDataArr.map((item: any) => {
            item.sub_nodes = []
            return item
        })

        fromDataArr = arrayToTree(fromDataArr)

        const removeEmptyDirectory = (arr: any, node: any, parent: any) => {
            let i = arr.length

            while (i--) {
                const item = arr[i]
                removeEmptyDirectory(item.sub_nodes, item, arr)
            }

            if (node && node.type === DOCUMENT_TYPES.DIR) {
                const child = node.sub_nodes || []
                const len = child.length
                if (!len) {
                    const idx = parent.indexOf(node)
                    idx != -1 && parent.splice(idx, 1)
                }
            }
        }

        // 移除所有空目录
        removeEmptyDirectory(fromDataArr, null, null)

        toDataArr.map((item: any) => {
            item.sub_nodes = []
            return item
        })

        fromData.value = fromDataArr
        toData.value = arrayToTree(toDataArr)
    })

    const sortTree = (nodes: any) => {
        ;(nodes || []).forEach((item: any) => {
            if (item.sub_nodes && item.sub_nodes.length) {
                sortTree(item.sub_nodes)
            }
        })
        nodes.sort((pre: any, next: any) => pre.sort - next.sort)
    }

    // 重新排序
    const onChangeTree = (f: any, d: any) => {
        const from = arrayToTree(flattenDeep(f || [], defaultProps.children))
        const to = arrayToTree(flattenDeep(d || [], defaultProps.children))

        sortTree(from)
        sortTree(to)

        fromData.value = from
        toData.value = to

        // setTimeout(() => {
        //     expandedKeys.value = toRaw(expandedKeys || [])
        // }, 2000)
    }

    defineExpose({
        show,
    })
</script>
