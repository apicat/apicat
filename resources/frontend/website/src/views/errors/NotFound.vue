<template>
    <div class="text-center">
        <img class="not-found-image" src="@/assets/image/404@2x.png" alt="404" />
        <p class="not-found-tip">啊哦，网页走丢了，正在努力寻找中…<router-link class="text-blue-600" :to="homePath">回到首页</router-link></p>
    </div>
</template>
<script setup lang="ts">
    import { inject } from 'vue'
    import { useUserStore } from '@/stores/user'
    import { ROOT_PATH, MAIN_PATH } from '@/router/constant'
    const switchNavState = inject('switchNavState') as any
    const homePath = { path: useUserStore().isLogin ? MAIN_PATH : ROOT_PATH }

    switchNavState && switchNavState(false)
</script>
<style>
    .not-found-image {
        height: 282px;
        margin-bottom: 6vh;
        margin-top: 20vh;
        display: inline-block;
    }

    .not-found-tip {
        font-size: 18px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 400;
        color: rgba(0, 0, 0, 0.65);
        line-height: 25px;
        letter-spacing: 4px;
    }
</style>
