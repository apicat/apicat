<template>
    <header class="h-14">
        <div class="shadow bg-white fixed top-0 w-full z-50">
            <div class="viewport relative flex items-center justify-between h-14">
                <div class="flex-1 flex items-center">
                    <router-link class="inline-flex items-center" to="/main">
                        <img class="h-10 w-auto" src="@/assets/image/logo.svg" alt="ApiCat" /><span class="logo-text logo-apicat">ApiCat</span>
                    </router-link>
                    <div class="ml-6">
                        <router-link v-for="item in navs" :to="{ name: item.name }" :key="item.name" :class="navClass(item)">{{ item.title }}</router-link>
                    </div>
                </div>
                <div class="flex items-center">
                    <div class="ml-3 relative">
                        <el-dropdown placement="bottom-end" trigger="click" popper-class="nav-dropdown" :popper-options="popperOptions">
                            <span class="el-dropdown-link">
                                <button type="button" class="flex items-center rounded-full" aria-expanded="false" aria-haspopup="true">
                                    <el-avatar>
                                        <img v-if="userInfo.avatar" :src="userInfo.avatar" />
                                        <span v-else>{{ lastName }}</span>
                                    </el-avatar>

                                    <el-icon class="el-icon--right"><arrow-down /></el-icon>
                                </button>
                            </span>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <div class="flex flex-col divide-y text-sm justify-center">
                                        <div class="px-6 py-2">
                                            <p class="text-base font-bold text-gray-600 truncate" :title="userInfo.name">{{ userInfo.name }}</p>
                                            <p class="text-gray-500 truncate" :title="userInfo.email">{{ userInfo.email }}</p>
                                        </div>
                                        <router-link
                                            v-for="item in dropMenusNavs"
                                            :key="item.name"
                                            class="flex items-center px-6 py-3 hover:text-gray-700"
                                            :to="{ name: item.name }"
                                        >
                                            <i :class="'iconfont ' + item.icon"></i><span class="ml-1">{{ item.title }}</span>
                                        </router-link>
                                        <a class="flex items-center px-6 py-3 hover:text-gray-700" href="javascript:void(0)" @click="onLogoutClick">
                                            <i class="iconfont iconIconPopoverExit"></i><span class="ml-1">退出登录</span>
                                        </a>
                                    </div>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </div>
                </div>
            </div>
        </div>
    </header>
</template>
<script lang="ts" setup>
    import { onMounted, ref } from 'vue'
    import { ArrowDown } from '@element-plus/icons-vue'
    import { useUserStore } from '@/stores/user'
    import { storeToRefs } from 'pinia'
    import { useRouter } from 'vue-router'

    const popperOptions = {
        modifiers: [
            {
                name: 'offset',
                options: {
                    offset: [-24, 10],
                },
            },
        ],
    }

    const dropMenusNavs = ref([
        { name: 'user.profile', title: '个人信息', icon: 'iconIconPopoverUser' },
        // { name: 'setting', title: '系统设置', icon: 'iconIconPopoverSetting' },
    ])

    const navs = [
        { name: 'projects', title: '项目', active: 'project' },
        { name: 'members', title: '成员' },
    ]

    const userStore = useUserStore()
    const { currentRoute } = useRouter()
    const { userInfo, lastName } = storeToRefs(userStore)

    const onLogoutClick = async () => {
        await userStore.logout()
    }

    const navClass = (item: any) => {
        return [
            'px-3',
            {
                'text-blue-600': ((currentRoute.value.name as string) || '').startsWith(item.active || item.name),
            },
        ]
    }

    onMounted(async () => {
        await userStore.getUserInfo()
    })
</script>
