<template>
    <el-popover v-model:visible="visible" ref="popoverRef" placement="bottom-end" :virtual-ref="searchDocumentPopoverRefEl" virtual-triggering width="auto">
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
    import { ref, unref } from 'vue'
    import { debounce } from 'lodash-es'
    import { useRoute, useRouter } from 'vue-router'
    import { searchDocuments } from '@/api/document'
    import { toDocumentDetailPath } from '@/router/document.router'
    import { traverseTree } from '@natosoft/shared'
    import { useDocumentStore } from '@/stores/document'
    import { storeToRefs } from 'pinia'

    const { params } = useRoute()
    const { push } = useRouter()
    const { apiDocTree } = storeToRefs(useDocumentStore())

    const visible = ref(false)
    const isSearched = ref(false)
    const isSearching = ref(false)
    const keywords = ref('')
    const searchData: any = ref([])
    const searchDomRef = ref()
    const popoverRef = ref()
    const searchInput = ref()
    const searchDocumentPopoverRefEl = ref()

    onClickOutside(
        searchDomRef,
        () => {
            visible.value = false
            resetSearch()
        },
        {
            ignore: [searchDocumentPopoverRefEl],
        }
    )

    const onSearchDocument = debounce(() => {
        const unRefKeywords = unref(keywords)

        if (unRefKeywords.length < 2) {
            return
        }

        isSearched.value = true
        isSearching.value = true

        searchDocuments(params.project_id, unRefKeywords)
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
                if (item.id === parseInt(doc.node_id, 10)) {
                    node = item
                    return false
                }
            },
            apiDocTree.value || [],
            { subKey: 'sub_nodes' }
        )

        const path = toDocumentDetailPath({ project_id: params.project_id, node_id: doc.node_id })

        if (node) {
            push({ path })
        } else {
            location.href = path
        }
    }

    const show = (el: any) => {
        searchDocumentPopoverRefEl.value = el
        visible.value = true
        searchInput.value.focus()
    }

    defineExpose({
        show,
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
