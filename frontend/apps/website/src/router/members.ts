import type { RouteRecordRaw } from 'vue-router'

const MemberListPage = () => import('@/views/member/MemberListPage.vue')

export const membersRoute: RouteRecordRaw = {
  name: 'members',
  path: '/members',
  component: MemberListPage,
  meta: {
    title: '成员列表',
  },
}
