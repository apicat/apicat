<template>
  <h4 v-if="readonly">{{ definition.description }}</h4>

  <template v-else>
    <input class="ac-document__title" type="text" v-input-limit v-model="definition.name" maxlength="255" ref="title" :placeholder="$t('app.schema.form.title')" />
    <input class="ac-document__desc" type="text" v-model="definition.description" maxlength="255" ref="title" :placeholder="$t('app.schema.form.desc')" />
  </template>

  <div class="ac-editor mt-10px">
    <JSONSchemaEditor :readonly="readonly" v-model="definition.schema" :definitions="definitions" />
  </div>
</template>
<script setup lang="ts">
import JSONSchemaEditor from '@/components/APIEditor/Editor.vue'
import { DefinitionSchema } from '@/components/APIEditor/types'

const props = defineProps<{
  modelValue: DefinitionSchema
  readonly?: boolean
  // 引用模型的集合
  definitions?: DefinitionSchema[]
}>()

const emit = defineEmits(['update:modelValue'])
const definition = useVModel(props, 'modelValue', emit)
</script>
