<template>
    <div class="loading">
        <el-icon class="is-loading" :size="30">
            <Loading />
        </el-icon>
    </div>
</template>
<script setup lang="ts">
    import { Loading } from '@element-plus/icons-vue'
</script>
<style lang="scss" scoped>
    .loading {
        position: absolute;
        top: 60px;
        bottom: 0;
        right: 0;
        left: 0;
        background: #fff;
        text-align: center;
        padding-top: 20vh;
    }
</style>
