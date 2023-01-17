<template>
    <el-popover
        ref="popoverRef"
        placement="bottom-end"
        trigger="click"
        :virtual-ref="props.virtualRef"
        virtual-triggering
        width="auto"
        :popper-options="popperOptions"
        @after-enter="onShow"
        @after-leave="onHide"
    >
        <div ref="searchDomRef" class="search-poptip">
            <el-input ref="searchInput" :suffix-icon="Search" placeholder="关键字 回车搜索" v-model="keywords" @change="onSearchDocument" clearable />
            <div v-show="isSearched" class="text-left search-result scroll-content" v-loading="isSearching">
                <span class="cursor-pointer" v-for="(item, index) in searchData" :key="item.doc_id + '_' + index" @click="onSearchResultItemClick(item)">
                    <h4 class="px-2 truncate hover:bg-gray-100">
                        {{ item.title }}
                    </h4>
                </span>

                <Result v-show="isSearched && !isSearching && !searchData.length">
                    <template #icon>
                        <img src="@/assets/image/icon-no-result.png" />
                    </template>
                    <template #title>
                        <span>未找到搜索结果！</span>
                    </template>
                </Result>
            </div>
        </div>
    </el-popover>
</template>
<script setup lang="ts">
    import { Search } from '@element-plus/icons-vue'
    import { onClickOutside } from '@vueuse/core'
    import { ref, unref, toRefs } from 'vue'
    import { debounce } from 'lodash-es'
    import { useRouter } from 'vue-router'
    import { searchDocuments, toDocumentDetailPath } from '@/api/document'
    import { toIterateDocumentDetailPath } from '@/api/iterate'
    import { traverseTree } from '@natosoft/shared'
    import { useDocumentStore } from '@/stores/document'
    import { storeToRefs } from 'pinia'
    import useIdPublicParam, { generateProjectOrIterateParams, getIdPublicByRouter } from '@/hooks/useIdPublicParam'

    const props = defineProps({
        virtualRef: Object,
    })

    const { push } = useRouter()
    const { apiDocTree } = storeToRefs(useDocumentStore())
    const publicParams = useIdPublicParam()
    const id_public: any = getIdPublicByRouter()
    const { isIterateRoute } = publicParams

    const isSearched = ref(false)
    const isSearching = ref(false)
    const keywords = ref('')
    const searchData: any = ref([])
    const searchDomRef = ref()
    const popoverRef = ref()
    const searchInput = ref()

    const propsRef: any = toRefs(props)

    const popperOptions = {
        modifiers: [
            { name: 'computeStyles', options: { gpuAcceleration: false } },
            { name: 'offset', options: { offset: [-300, 10] } },
        ],
    }

    onClickOutside(searchDomRef, () => resetSearch(), {
        ignore: [propsRef.virtualRef],
    })

    const onSearchDocument = debounce(() => {
        const unRefKeywords = unref(keywords)

        if (unRefKeywords.length < 2) {
            return
        }

        isSearched.value = true
        isSearching.value = true
        const data = generateProjectOrIterateParams(publicParams)
        data.project_id_public = publicParams.projectPublicId

        searchDocuments(data, unRefKeywords)
            .then((res) => {
                searchData.value = res.data || []
            })
            .catch((e) => {
                //
            })
            .finally(() => {
                isSearching.value = false
            })
    }, 200)

    const resetSearch = () => {
        searchData.value = []
        keywords.value = ''
        isSearched.value = false
    }

    const onSearchResultItemClick = (doc: any) => {
        // 检测是否存在该文档，存在内部跳转。不存在 href
        let node = null

        traverseTree(
            (item: any) => {
                if (item.id === parseInt(doc.doc_id, 10)) {
                    node = item
                    return false
                }
            },
            apiDocTree.value || [],
            { subKey: 'sub_nodes' }
        )

        const path = isIterateRoute ? toIterateDocumentDetailPath(id_public, doc.doc_id) : toDocumentDetailPath(id_public, doc.doc_id)

        if (node) {
            push({ path })
        } else {
            location.href = path
        }
    }

    const onShow = () => {
        searchInput.value.focus()
    }

    const onHide = () => {
        resetSearch()
    }

    defineExpose({
        onShow,
        onHide,
    })
</script>
<style lang="scss" scoped>
    .search-poptip {
        width: 300px;
        margin: -12px;
        padding: 14px 24px;
        .search-result {
            max-height: 220px;
            position: relative;
            line-height: 28px;
            overflow: scroll;
            padding-top: 10px;

            h4 {
                cursor: pointer;
            }
        }
    }
</style>
