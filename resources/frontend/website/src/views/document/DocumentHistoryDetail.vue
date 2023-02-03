<template>
    <div class="ac-document document-history-detail" v-loading="isLoading" element-loading-background="#fff">
        <div v-show="hasDocument && document.id">
            <h1 class="ac-document__title" ref="title">{{ document.title }}</h1>
            <div v-if="document.content" class="ProseMirror readonly" v-html="document.content" />
        </div>

        <div v-if="!hasDocument">
            <Result :styles="{ width: '260px', height: 'auto', 'margin-bottom': '26px' }">
                <template #icon>
                    <img src="@/assets/image/icon-empty.png" alt="" />
                </template>
                <template #title>
                    <div style="width: 470px; display: block; margin: auto">暂无文档历史记录</div>
                </template>
            </Result>
        </div>

        <ac-backtop :bottom="100" :right="100" />
    </div>
</template>
<script lang="ts" setup>
    import { getDocumentHistoryRecordDetail } from '@/api/document'
    import { onMounted, onUnmounted, ref, watch } from 'vue'
    import { hideLoading } from '@/hooks/useLoading'
    import { useDocumentDetailInteractive } from '@/hooks/useDocumentDetailInteractive'
    import { useRoute } from 'vue-router'

    const route = useRoute()
    const isLoading = ref(false)
    const hasDocument = ref(true)
    const document: any = ref({})

    const { project_id } = route.params

    watch(
        () => route.params.id,
        () => getDocumentDetail()
    )

    const getDocumentDetail = async () => {
        const id = parseInt(route.params.id as string, 10)

        if (isNaN(id)) {
            isLoading.value = false
            hasDocument.value = false
            return
        }

        isLoading.value = true
        hasDocument.value = true

        getDocumentHistoryRecordDetail(project_id, id)
            .then((res) => {
                document.value = res.data || {}
                useDocumentDetailInteractive('.document-history-detail')
            })
            .catch((e) => {
                //
            })
            .finally(() => {
                isLoading.value = false
                setTimeout(() => hideLoading(), 500)
            })
    }

    onMounted(async () => {
        await getDocumentDetail()
    })

    onUnmounted(() => {
        isLoading.value = true
        hasDocument.value = true
        document.value = {}
    })
</script>
