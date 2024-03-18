<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useNamespace } from '@apicat/hooks'
import { useI18n } from 'vue-i18n'
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import { useUserStore } from '@/store/user'
import { USER_PAGE_NAME } from '@/router/constant'
import { useTitle } from '@/hooks/useTitle'

const router = useRouter()
const route = useRoute()

const ns = useNamespace('group-list')
const { t } = useI18n()
const browserTitle = useTitle()

const menus: Record<
  string,
  {
    icon: string
    title: string
    component: any
  }
> = {
  general: {
    icon: 'uil:user',
    component: defineAsyncComponent(() => import('./pages/General.vue')),
    title: t('app.pageTitles.userSetting.general'),
  },
  email: {
    icon: 'ic:outline-email',
    component: defineAsyncComponent(() => import('./pages/Email.vue')),
    title: t('app.pageTitles.userSetting.email'),
  },
  github: {
    icon: 'mdi:github',
    component: defineAsyncComponent(() => import('./pages/Github.vue')),
    title: t('app.pageTitles.userSetting.github'),
  },
  password: {
    icon: 'iconamoon:shield-yes-light',
    component: defineAsyncComponent(() => import('./pages/Password.vue')),
    title: t('app.pageTitles.userSetting.password'),
  },
}

const currentPage = computed<string>(() => (route.params.page as string) || Object.keys(menus)[0])

// set default browser title
browserTitle.value = menus[currentPage.value].title

function setCurrentPage(menuKey: string) {
  router.push({
    name: USER_PAGE_NAME,
    params: {
      page: menuKey,
    },
  }).then(() => {
    browserTitle.value = menus[menuKey].title
  })
}

function activeClass(key: string) {
  return currentPage.value === key ? 'active' : ''
}

const userStore = useUserStore()
userStore.getUserInfo()
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
            <Icon
              color="rgba(72,148,255)"
              class="mr-2"
              :icon="menu.icon"
              width="18"
            />
            <span>{{ t(`app.user.${key}.left_title`) }}</span>
          </div>
        </div>
      </div>
    </template>
    <component :is="menus[currentPage].component" />
  </LeftRightLayout>
</template>
