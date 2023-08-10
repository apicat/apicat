<template>
  <div class="flex h-full">
    <div class="w-246px b-r b-solid b-gray-110 px-20px">
      <IterationTree v-model:selected-project-key="selectedProjectKeyRef" :projects="followedProjects" />
    </div>
    <div class="flex-1">
      <IterationTable v-show="isListMode" :data="iterations" v-model:page="currentPage" :total="total" />
      <IterationForm v-show="isFormMode" />
    </div>
  </div>
</template>
<script setup lang="ts">
import IterationTree from './components/IterationTree.vue'
import IterationTable from './components/IterationTable.vue'
import IterationForm from './components/IterationForm.vue'
import { usePageMode } from './logic/usePageMode'
import { useFollowedProjectList } from './logic/useFollowedProjectList'
import { useIterationList } from './logic/useIterationList'

// import components
// IterationTree
// props[IterationTreeData,activeKey:string];
// event[add,click]);

// IterationTable
// (props[IterationTableData,pageSize:number,total:number];
// v-model:page;
// event[delete,edit];

// IterationForm
// props[iteration|null,handleSubmit:Promise<void>];
// event[add(iteration),update(iteration)]);

const { isFormMode, isListMode, toggleMode } = usePageMode()
const { selectedProjectKeyRef, followedProjects, setSelectedProjectKey } = useFollowedProjectList()
const { data: iterations, currentPage, total, fetchIterationList } = useIterationList(selectedProjectKeyRef)

// todo 获取所有迭代目录树

// todo 根据当前所选迭代获取该迭代下的迭代列表

// todo 添加迭代

// todo 编辑迭代
</script>
