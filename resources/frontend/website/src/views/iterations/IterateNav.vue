<template>
    <SidebarLayout>
        <div class="flex flex-col border divide-y bg-white">
            <a
                class="relative h-12 pl-6 flex items-center text-neutral-600 hover:text-neutral-900"
                href="javascript:void(0);"
                :class="setActiveClass({ project_id: '' })"
                @click="onCollectItemClick({})"
            >
                <span class="ml-1">迭代列表</span>
            </a>
        </div>

        <div class="overflow-y-scroll mt-2 collection_list" v-loading="isLoading">
            <div class="flex flex-col border divide-y bg-white" ref="sortableEle">
                <a
                    v-for="i in collection"
                    :key="i.project_id"
                    class="relative flex pr-4 items-center justify-between text-neutral-600 hover:text-neutral-900 collection_item"
                    :class="setActiveClass(i)"
                    href="javascript:void(0);"
                    @click="onCollectItemClick(i)"
                >
                    <section class="overflow-hidden flex items-center flex-1 h-12 pl-6">
                        <span class="ml-1 truncate" :title="i.project_name">{{ i.project_name }}</span>
                    </section>
                </a>

                <a class="relative flex pr-4 items-center justify-between text-neutral-600 hover:text-neutral-900" v-if="!collection.length">
                    <section class="overflow-hidden flex items-center flex-1 h-12 pl-6">
                        <span class="ml-1" title="ApiCat">暂无收藏</span>
                    </section>
                </a>
            </div>
        </div>
    </SidebarLayout>
</template>

<script setup lang="ts">
    import SidebarLayout from '@/layout/SidebarLayout.vue'
    import { onMounted, onUnmounted, ref } from 'vue'
    import { useIterateStore } from '@/stores/iterate'
    import { useSortable } from '@/hooks/useSortable'
    import { storeToRefs } from 'pinia'
    import useApi from '@/hooks/useApi'

    const sortableEle = ref()
    const iterateStore = useIterateStore()
    const { activeTab, collection } = storeToRefs(iterateStore)

    const [isLoading, execute] = useApi(iterateStore.getIterateCollectionList)

    const setActiveClass = (item: any) => [{ active: activeTab.value === item.project_id }]

    const onCollectItemClick = (item: any) => {
        iterateStore.switchActiveCollectTab(item.project_id)
    }

    const { initSortable } = useSortable(sortableEle, {
        draggable: '.collection_item',
        onEnd(e: any) {
            const { oldIndex, newIndex } = e
            iterateStore.sortCollectionList(oldIndex, newIndex)
        },
    })

    onMounted(async () => {
        await execute()
        initSortable()
    })

    onUnmounted(() => {
        iterateStore.$patch({ activeTab: undefined })
    })
</script>

<style lang="scss" scoped>
    .collection_list {
        max-height: 491px;

        .collection_item {
            &.sortable-drag {
                background-color: rgb(250, 250, 250);
            }
        }
    }
</style>
