<template>
    <div class="ac-header-operate">
        <div class="ac-header-operate__main">
            <p :class="titleClass">13123</p>
        </div>

        <div class="ac-header-operate__btns">
            <el-button type="primary">保存</el-button>
            <i class="iconfont iconshare2"></i>
            <i class="iconfont iconIconPopoverUpload"></i>
        </div>
    </div>
</template>
<script setup lang="ts">
    import { ref, computed } from 'vue'
    import emitter, { IS_SHOW_DOCUMENT_TITLE } from '@/common/emitter'

    const isShowTitle = ref(false)
    const titleClass = computed(() => {
        return [
            'ac-header-operate__title animate__animated animate__faster',
            {
                hidden: !isShowTitle.value,
                animate__slideInDown: isShowTitle.value,
                animate__slideOutUp: !isShowTitle.value,
            },
        ]
    })

    emitter.on(IS_SHOW_DOCUMENT_TITLE, (isShow: any) => {
        isShowTitle.value = isShow
    })
</script>
