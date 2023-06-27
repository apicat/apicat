<template>
  <div class="flex gap-30px" :class="{ 'justify-around': exportType === 'collection' }">
    <div class="flex flex-col items-center w-100px" v-for="item in exportList" :key="item.text" @click="handleSelect(selectedRef, item)">
      <div class="border border-solid rounded cursor-pointer p-20px hover:border-blue-primary" :class="[{ 'border-blue-primary': selectedRef.type === item.type }]">
        <img :src="item.logo" alt="" class="w-60px h-60px" />
      </div>
      <p>{{ item.text }}</p>
      <div v-if="item.versions" class="mt-4px">
        <el-select v-model="selectedRef.version" size="small" @change="handleSelect(selectedRef, item)">
          <el-option v-for="v in item.versions" :label="v.label" :value="v.value"> </el-option>
        </el-select>
      </div>
    </div>
  </div>

  <div :class="{ 'text-right': exportType === 'collection' }">
    <el-button type="primary" class="mt-20px" @click="handleExport(selectedRef)">{{ $t('app.common.export') }}</el-button>
  </div>
</template>
<script setup lang="ts">
import swaggerLogo from '@/assets/images/logo-swagger@2x.png'
import openApiLogo from '@/assets/images/logo-openapis.svg'
import htmlLogo from '@/assets/images/logo-html@2x.png'
import mdLogo from '@/assets/images/logo-markdown@2x.png'
import { ExportProjectTypes } from '@/commons/constant'
import { exportProject } from '@/api/project'
import { useParams } from '@/hooks/useParams'
import { exportCollection } from '@/api/collection'
interface ExportParams {
  type: ExportProjectTypes
  version: string
  params: Record<string, any>
}

interface ExportItem {
  logo: string
  text: string
  type: ExportProjectTypes
  params?: Record<string, any>
  versions?: Array<{
    label: string
    value: string
  }>
}

const props = withDefaults(
  defineProps<{
    exportType?: 'project' | 'collection'
    project_id?: string
    doc_id?: string
  }>(),
  {
    exportType: 'project',
  }
)

let { project_id, doc_id } = useParams()

const exportList: ExportItem[] = [
  { logo: swaggerLogo, text: 'Swagger 2.0', type: ExportProjectTypes.Swagger },
  {
    logo: openApiLogo,
    text: 'OpenAPI',
    type: ExportProjectTypes.OpenAPI,
    versions: [
      { label: '3.0.0', value: 'openapi3.0.0' },
      { label: '3.0.1', value: 'openapi3.0.1' },
      { label: '3.0.2', value: 'openapi3.0.2' },
      { label: '3.1.0', value: 'openapi3.1.0' },
    ],
  },
  { logo: htmlLogo, text: 'HTML', type: ExportProjectTypes.HTML, params: { download: true } },
  { logo: mdLogo, text: 'Markdown', type: ExportProjectTypes.MARKDOWN, params: { download: true } },
]

const selectedRef: Ref<ExportParams> = ref({
  type: exportList[0].type,
  version: 'openapi3.0.0',
  params: {},
})

const handleSelect = (selected: ExportParams, item: any) => {
  selected.type = item.type
  selected.params = item.params || {}
}

const handleExport = (selected: ExportParams) => {
  let type: string = selected.type

  if (selected.type === ExportProjectTypes.OpenAPI) {
    type = selected.version
  }

  project_id = props.project_id || (project_id as string)
  doc_id = props.doc_id || (doc_id as string)

  const params = { type, ...selected.params }

  const exportUrl = props.exportType === 'project' ? exportProject({ project_id, ...params }) : exportCollection({ project_id, collection_id: doc_id, ...params })
  window.open(exportUrl)
}
</script>
