<template>
  <div class="container flex flex-col justify-center mx-auto" v-loading="isLoading">
    <p class="border-b border-solid border-gray-lighter py-20px mt-20px text-20px">
      {{ $t('app.project.list.tabTitle') }}
    </p>

    <ul class="ac-project-list my-20px py-10px" v-if="projects.length">
      <li
        class="flex flex-col overflow-hidden rounded shadow-md cursor-pointer hover:shadow-lg w-250px h-156px"
        v-for="project in projects"
        @click="$router.push(getProjectDetailPath(project.id))"
      >
        <div class="flex items-center justify-center h-112px bg-#6699CC">
          <Iconfont class="text-white" icon="ac-xiangmu" :size="55" />
        </div>
        <p class="flex items-center flex-1 truncate px-16px">{{ project.title }}</p>
      </li>

      <li class="flex flex-col justify-between rounded cursor-pointer w-250px h-156px hover:shadow-lg bg-gray-110 px-20px py-16px" @click="handleShowModelClick">
        <ac-icon-ep-plus class="text-18px" />
        <p>{{ $t('app.project.createModal.title') }}</p>
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

onBeforeMount(async () => await getProjectListApi())
onBeforeUnmount(() => (projectStore.projects = []))
</script>
<style scoped>
.ac-project-list {
  display: grid;
  justify-content: space-between;
  grid-template-columns: repeat(auto-fill, 250px);
  grid-gap: 20px;
}
</style>
