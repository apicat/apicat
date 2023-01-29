<template>
    <header class="h-14">
        <div class="shadow bg-white fixed top-0 w-full z-50">
            <div class="viewport relative items-center justify-between h-14 lg:flex hidden">
                <div class="flex-1 flex items-center">
                    <div class="flex-shrink-0 flex">
                        <a class="inline-flex items-center" href="/">
                            <img class="h-10 w-auto" src="@/assets/image/logo.svg" alt="ApiCat" /><span class="logo-text logo-apicat">ApiCat</span>
                        </a>
                    </div>
                    <div class="ml-6 hidden lg:block" v-if="isShowNav">
                        <div class="flex text-gray-600">
                            <a href="https://apicat.net" target="_blank" class="px-3 py-2 rounded-md text-sm">首页</a>
                            <a href="https://apicat.net/app/31eb8ec8c295c8c4130ee50304d55ab5" target="_blank" class="px-3 py-2 rounded-md text-sm">文档</a>
                            <a href="https://github.com/NatosoftCN/ApiCat_Open" target="_blank" class="px-3 py-2 rounded-md text-sm">Github</a>
                        </div>
                    </div>
                </div>

                <div v-if="isShowNav">
                    <router-link :to="LOGIN_PATH" class="mr-4"><el-button>登录</el-button></router-link>
                    <router-link :to="REGISTE_PATH"><el-button type="primary">注册</el-button></router-link>
                </div>
            </div>

            <div class="lg:hidden">
                <div class="flex items-center justify-between h-14 px-4 relative">
                    <router-link class="inline-flex items-center" to="/">
                        <img class="h-10 w-auto" src="@/assets/image/logo.svg" alt="ApiCat" /><span class="logo-text logo-apicat">ApiCat</span>
                    </router-link>

                    <div class="cursor-pointer" @click="isShow = !isShow">
                        <svg viewBox="0 0 24 24" width="1em" height="1em" style="font-size: 24px" fill="currentColor">
                            <path d="M21 7H3V5H21V7Z" fill-rule="evenodd" clip-rule="evenodd" fill="currentColor"></path>
                            <path d="M21 13H3V11H21V13Z" fill-rule="evenodd" clip-rule="evenodd" fill="currentColor"></path>
                            <path d="M21 19H3V17H21V19Z" fill-rule="evenodd" clip-rule="evenodd" fill="currentColor"></path>
                        </svg>
                    </div>

                    <Transition name="slide-fade">
                        <div class="bg-white fixed top-0 left-0 w-full max-h-screen flex flex-col z-10" v-if="isShow">
                            <div class="flex items-center justify-between h-14 px-4">
                                <router-link class="inline-flex items-center" to="/">
                                    <img class="h-10 w-auto" src="@/assets/image/logo.svg" alt="ApiCat" /><span class="logo-text logo-apicat">ApiCat</span>
                                </router-link>

                                <div class="cursor-pointer" @click="isShow = !isShow">
                                    <svg viewBox="0 0 24 24" width="1em" height="1em" style="font-size: 24px" fill="currentColor">
                                        <path
                                            d="M12.0001 13.4141L17.6568 19.0708L19.0711 17.6566L13.4143 11.9998L19.0711 6.34304L17.6569 4.92883L12.0001 10.5856L6.34314 4.92871L4.92893 6.34292L10.5859 11.9998L4.92896 17.6567L6.34317 19.071L12.0001 13.4141Z"
                                            fill-rule="evenodd"
                                            clip-rule="evenodd"
                                            fill="currentColor"
                                        ></path>
                                    </svg>
                                </div>
                            </div>
                            <div class="flex-1 overflow-y-auto">
                                <div class="px-4 text-base font-medium">
                                    <a class="px-3 py-2 mb-1 rounded-lg block" href="/">首页</a>
                                    <a class="px-3 py-2 mb-1 rounded-lg block" href="#">文档</a>
                                    <a class="px-3 py-2 mb-1 rounded-lg block" href="#">GitHub</a>
                                </div>
                                <div class="px-5">
                                    <router-link to="/login"><el-button class="w-full">登录</el-button></router-link>
                                    <router-link to="/"><el-button type="primary" class="w-full my-5">注册</el-button></router-link>
                                </div>
                            </div>
                        </div>
                    </Transition>
                    <div class="fixed top-0 left-0 w-screen h-screen bg-opacity-40 bg-gray-900" v-if="isShow" @click="isShow = !isShow"></div>
                </div>
            </div>
        </div>
    </header>
</template>

<script setup lang="ts">
    import { LOGIN_PATH, REGISTE_PATH } from '@/router/constant'
    import { ref, watch } from 'vue'
    withDefaults(defineProps<{ isShowNav: boolean }>(), {
        isShowNav: true,
    })

    const isShow = ref(false)

    watch(isShow, () => {
        document.body.classList.toggle('overflow-hidden')
    })
</script>
<style scoped>
    .slide-fade-enter-active,
    .slide-fade-leave-active {
        transition: all 0.3s ease-out;
    }

    .slide-fade-enter-from,
    .slide-fade-leave-to {
        transform: translateY(-100%);
    }
</style>
