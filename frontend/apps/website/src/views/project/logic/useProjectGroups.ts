// import { storeToRefs } from 'pinia'
// import { useI18n } from 'vue-i18n'
// import type CreateOrUpdateProjectGroup from '../components/CreateOrUpdateProjectGroup.vue'
// import { sortProjectGroup } from '@/api/projectGroup'
// import { AsyncMsgBox } from '@/components/AsyncMessageBox'
// import useApi from '@/hooks/useApi'
// import useProjectGroupStore from '@/store/projectGroup'
// import type { ProjectGroup, ProjectGroupSelectKey } from '@/typings'

// export function useProjectGroups() {
//   const { t } = useI18n()
//   const createOrUpdateProjectGroupRef = ref<InstanceType<typeof CreateOrUpdateProjectGroup>>()
//   const selectedGroupRef = ref<ProjectGroupSelectKey>('all')
//   const projectGroupStore = useProjectGroupStore()
//   const { projectGroups, groupsForOptions } = storeToRefs(projectGroupStore)
//
//   const [isLoading, getProjectGroupsApi] = useApi(projectGroupStore.getProjectGroups)
//
//   // 删除项目分组
//   const handleDeleteProjectGroup = async (group: ProjectGroup) => {
//     AsyncMsgBox({
//       title: t('app.common.deleteTip'),
//       content: '确定删除该项目分组吗?',
//       onOk: async () => {
//         // 删除选中分组
//         if (selectedGroupRef.value === group.id)
//           selectedGroupRef.value = 'all'
//
//         await projectGroupStore.deleteProjectGroup(group.id!)
//         await getProjectGroupsApi()
//       },
//     })
//   }
//
//   // 重命名项目分组
//   const handleRenameProjectGroup = (group?: ProjectGroup) => createOrUpdateProjectGroupRef.value?.show(group)
//
//   // 项目分组排序
//   const handleSortProjectGroup = async () => await sortProjectGroup(projectGroups.value.map(item => item.id!))
//
//   // 获取项目分组
//   onMounted(async () => await getProjectGroupsApi())
//
//   return {
//     createOrUpdateProjectGroupRef,
//     selectedGroupRef,
//     isLoading,
//     projectGroups,
//     groupsForOptions,
//
//     refreshProjectGroups: getProjectGroupsApi,
//     handleDeleteProjectGroup,
//     handleRenameProjectGroup,
//     handleSortProjectGroup,
//   }
// }
