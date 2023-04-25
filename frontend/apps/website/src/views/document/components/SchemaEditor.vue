<template>
  <div>
    <div :class="ns.b()">
      <template v-if="readonly">
        <h4 class="">{{ definition.description }}</h4>
      </template>

      <template v-else>
        <input class="ac-document__title" type="text" v-input-limit v-model="definition.name" maxlength="255" ref="title" placeholder="请输入模型标题" />
        <input class="ac-document__desc" type="text" v-model="definition.description" maxlength="255" ref="title" placeholder="请输入模型描述" />
      </template>

      <div class="ac-editor mt-10px">
        <JSONSchemaEditor :readonly="readonly" v-model="definition.schema" :definitions="definitions" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks/useNamespace'
import JSONSchemaEditor from '@/components/APIEditor/Editor.vue'
import { Definition } from '@/components/APIEditor/types'

const ns = useNamespace('document')
const props = defineProps<{
  modelValue: Definition
  readonly?: boolean
  // 引用模型的集合
  definitions?: Definition[]
}>()

const emit = defineEmits(['update:modelValue'])
const definition = useVModel(props, 'modelValue', emit)
</script>
