<template>
    <SidebarLayout v-if="!isReader && projectInfo">
        <div class="flex flex-col bg-white border divide-y adwa">
            <router-link
                v-for="item in navs"
                :key="item.name"
                :to="{ name: item.name }"
                class="relative flex items-center h-12 pl-6 text-neutral-600 hover:text-neutral-900"
            >
                <span>{{ item.meta.title }}</span>
            </router-link>

            <router-link :to="documentDetailPath" class="relative flex items-center h-12 pl-6 text-blue-600">
                <span class="truncate">前往 {{ projectInfo.name }} 文档</span>
            </router-link>
        </div>
    </SidebarLayout>

    <!-- 普通成员 -->
    <RouterView v-if="isReader" />
</template>

<script setup lang="ts">
    import { useRoute } from 'vue-router'
    import { ref } from 'vue'

    import SidebarLayout from '@/layout/SidebarLayout.vue'
    import { ProjectRoutes } from '@/router/project.router'
    import { useProjectStore } from '@/stores/project'
    import { PROJECT_ROLES_KEYS } from '@/common/constant'

    import { toDocumentDetailPath } from '@/router/document.router'

    const route = useRoute()
    const navs: any = ref([])
    const { isReader, isDeveloper, projectInfo } = useProjectStore()

    const documentDetailPath = ref({ path: toDocumentDetailPath({ project_id: route.params.project_id }) })

    if (projectInfo) {
        //  all navs
        navs.value = ProjectRoutes
        // developer navs
        if (isDeveloper) {
            navs.value = ProjectRoutes.filter((item) => item.meta.role.indexOf(PROJECT_ROLES_KEYS.DEVELOPER) !== -1)
        }
    }
</script>
