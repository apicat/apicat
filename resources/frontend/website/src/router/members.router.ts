const MemberList = () => import('../views/members/MemberList.vue')

export const MembersRootRoute = {
    path: '/members',
    name: 'members',
    meta: { title: '成员' },
}

export default {
    ...MembersRootRoute,
    component: MemberList,
}
