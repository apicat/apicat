<template>
  <div class="container flex flex-col justify-center mx-auto">
    <p class="border-b border-solid border-slate-500 py-20px mt-20px text-20px">
      {{ $t('app.project.list.tabTitle') }}
    </p>

    <ul class="ac-project-list mt-20px" v-loading="isLoading">
      <li
        class="flex flex-col overflow-hidden rounded-lg shadow-md cursor-pointer hover:shadow-lg w-250px h-156px"
        v-for="project in projects"
        @click="$router.push(getProjectDetailPath(project.id))"
      >
        <div class="flex items-center justify-center h-112px bg-#6699CC">
          <Iconfont class="text-white" icon="ac-xiangmu" :size="55" />
        </div>
        <p class="flex items-center flex-1 truncate px-16px">{{ project.title }}</p>
      </li>

      <li class="w-250px h-156px rounded-lg cursor-pointer hover:shadow-lg bg-#F2F2F2 flex flex-col justify-between px-20px py-16px" @click="handleShowModelClick">
        <ac-icon-ep-plus class="text-18px" />
        <p>创建项目</p>
      </li>
    </ul>
  </div>

  <CreateProjectModal ref="createProjectModal" />
</template>

<script lang="ts" setup>
import { getProjectDetailPath } from '@/router/document'
import CreateProjectModal from './CreateProjectModal.vue'
import uesProjectStore from '@/store/project'
import useApi from '@/hooks/useApi'
import { storeToRefs } from 'pinia'

const createProjectModal = ref<InstanceType<typeof CreateProjectModal>>()
const handleShowModelClick = () => createProjectModal.value!.show()
const projectStore = uesProjectStore()
const { projects } = storeToRefs(projectStore)
const [isLoading, getProjectListApi] = useApi(projectStore.getProjects)()

onMounted(async () => await getProjectListApi())
</script>
<style>
.ac-project-list {
  display: grid;
  justify-content: space-between;
  grid-template-columns: repeat(auto-fill, 250px);
  grid-gap: 20px;
}
</style>