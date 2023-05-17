<template>
  <div :class="ns.b()">
    <el-input :class="ns.e('filter')" v-model="searchWord" clearable placeholder="response name" />
    <div :class="ns.e('content')">
      <ul v-show="filterableResponses.length">
        <li :class="ns.e('item')" v-for="response in filterableResponses" :key="response.id" @click="handleSelect(response)">
          <Iconfont icon="ac-response" />
          {{ response.name }}
        </li>
      </ul>
      <div :class="ns.e('empty')" v-show="!filterableResponses.length">
        <el-empty description="No Response List" :image-size="50" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { DefinitionResponse } from '@/typings'

const emits = defineEmits(['select'])
const props = withDefaults(
  defineProps<{
    responses: DefinitionResponse[]
  }>(),
  {
    responses: () => [],
  }
)

const { responses } = toRefs(props)
const ns = useNamespace('def-response-search-list')
const searchWord = ref('')
const filterableResponses = computed(() => responses.value.filter((item) => item.name.toLowerCase().includes(unref(searchWord).toLowerCase())))

const handleSelect = (response: DefinitionResponse) => emits('select', toRaw(response))
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@include b(def-response-search-list) {
  @apply px-10px py-10px;
  @include e(filter) {
  }

  @include e(content) {
    @apply max-h-260px mt-10px;
  }

  @include e(item) {
    @apply cursor-pointer hover:bg-gray-100 px-10px py-4px;

    @include when(active) {
      @apply text-blue-primary;
    }
  }
}
</style>
