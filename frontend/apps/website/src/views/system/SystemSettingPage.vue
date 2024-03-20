<script setup lang="ts">
import { useNamespace } from '@apicat/hooks'
import { useI18n } from 'vue-i18n'
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import { useUserStore } from '@/store/user'
import { SYSTEM_PAGE_NAME } from '@/router/constant'
import { useTitle } from '@/hooks/useTitle'
import Iconfont from '@/components/Iconfont.vue'
import Page30pLayout from '@/layouts/Page30pLayout.vue'

const router = useRouter()
const route = useRoute()

const ns = useNamespace('group-list')
const { t } = useI18n()
const browserTitle = useTitle()
const userStore = useUserStore()

const menus: Record<string, { icon: string; title: string; component: any }> = {
  service: {
    icon: 'ac-webpage',
    component: defineAsyncComponent(() => import('./pages/Service.vue')),
    title: t('app.pageTitles.systemSetting.service'),
  },
  oauth: {
    icon: 'ac-user-outline',
    component: defineAsyncComponent(() => import('./pages/OAuth.vue')),
    title: t('app.pageTitles.systemSetting.oauth'),
  },
  storage: {
    icon: 'ac-memory-one',
    component: defineAsyncComponent(() => import('./pages/Storage.vue')),
    title: t('app.pageTitles.systemSetting.storage'),
  },
  cache: {
    icon: 'ac-storage-card-one',
    component: defineAsyncComponent(() => import('./pages/Cache.vue')),
    title: t('app.pageTitles.systemSetting.cache'),
  },
  db: {
    icon: 'ac-data',
    component: defineAsyncComponent(() => import('./pages/Database.vue')),
    title: t('app.pageTitles.systemSetting.database'),
  },
  email: {
    icon: 'ac-mail',
    component: defineAsyncComponent(() => import('./pages/Email.vue')),
    title: t('app.pageTitles.systemSetting.email'),
  },
  model: {
    icon: 'ac-zhinengyouhua',
    component: defineAsyncComponent(() => import('./pages/Model.vue')),
    title: t('app.pageTitles.systemSetting.model'),
  },
  users: {
    icon: 'ac-people-outline',
    component: defineAsyncComponent(() => import('./pages/Users.vue')),
    title: t('app.pageTitles.systemSetting.users'),
  },
}

const currentPage = computed<string>(() => (route.params.page as string) || Object.keys(menus)[0])

function setCurrentPage(menuKey: string) {
  router
    .push({
      name: SYSTEM_PAGE_NAME,
      params: {
        page: menuKey,
      },
    })
    .then(() => {
      browserTitle.value = menus[menuKey].title
    })
}

function activeClass(key: string) {
  return currentPage.value === key ? 'active' : ''
}

onBeforeMount(() => {
  if (!Object.keys(menus).includes(route.params.page as string))
    setCurrentPage(Object.keys(menus)[0])

  // set default browser title
  browserTitle.value = menus[currentPage.value].title

  // get user info
  userStore.getUserInfo()
})
</script>

<template>
  <LeftRightLayout main-width="auto">
    <template #left>
      <div class="flex flex-col h-full">
        <div :class="[ns.b(), ns.m('header')]">
          <div
            v-for="(menu, key) in menus"
            :key="key"
            :class="[ns.e('item'), activeClass(key)]"
            @click="setCurrentPage(key)"
          >
            <Iconfont color="rgba(72,148,255)" class="mr-2" :icon="menu.icon" width="18" />
            <p class="truncate" :title="t(`app.system.${key}.left_title`)">
              {{ t(`app.system.${key}.left_title`) }}
            </p>
          </div>
        </div>
      </div>
    </template>
    <Page30pLayout>
      <component :is="menus[currentPage].component" />
    </Page30pLayout>
  </LeftRightLayout>
</template>
