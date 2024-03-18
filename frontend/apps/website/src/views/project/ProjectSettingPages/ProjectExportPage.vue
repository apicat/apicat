<script setup lang="ts">
import swaggerLogo from '@/assets/images/logo-swagger@2x.png'
import openApiLogo from '@/assets/images/logo-openapis.svg'
import htmlLogo from '@/assets/images/logo-html@2x.png'
import mdLogo from '@/assets/images/logo-markdown@2x.png'
import apiCatLogo from '@/assets/images/logo-square.svg'
import { ExportProjectTypes } from '@/commons/constant'
import { apiExportProject } from '@/api/project'
import { useParams } from '@/hooks/useParams'
import { apiExportCollection } from '@/api/project/collection'

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
    doc_id?: string | number
  }>(),
  {
    exportType: 'project',
  },
)

let { project_id } = useParams()
const { currentRoute } = useRouter()
const exportList: ExportItem[] = [
  { logo: apiCatLogo, text: 'ApiCat', type: ExportProjectTypes.ApiCat },
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
  {
    logo: htmlLogo,
    text: 'HTML',
    type: ExportProjectTypes.HTML,
    params: { download: true },
  },
  {
    logo: mdLogo,
    text: 'Markdown',
    type: ExportProjectTypes.MARKDOWN,
    params: { download: true },
  },
]

const selectedRef: Ref<ExportParams> = ref({
  type: exportList[0].type,
  version: '',
  params: {},
})

function handleSelect(selected: ExportParams, item: any) {
  if (selected.type === item.type)
    return

  selected.type = item.type
  selected.version = item.versions ? item.versions[0].value : ''
  selected.params = item.params || {}
}

async function handleExport(selected: ExportParams) {
  let type: string = selected.type

  if (selected.type === ExportProjectTypes.OpenAPI)
    type = selected.version

  project_id = props.project_id || (project_id as string)
  const collection_id = props.doc_id || (currentRoute.value.params.collectionID as string)

  const exportUrl
    = props.exportType === 'project'
      ? (await apiExportProject(project_id, type, selected.params.download)).path
      : (await apiExportCollection(project_id, collection_id as number, type, selected.params.download)).path

  if (exportUrl)
    window.open(exportUrl)
}
</script>

<template>
  <div class="flex items-start gap-30px min-h-152px" :class="{ 'justify-around': exportType === 'collection' }">
    <div
      v-for="item in exportList"
      :key="item.text"
      class="flex flex-col items-center cursor-pointer w-100px "
      @click="handleSelect(selectedRef, item)"
    >
      <div
        class="border border-solid rounded p-20px hover:border-blue-primary"
        :class="[{ 'border-blue-primary': selectedRef.type === item.type }]"
      >
        <img :src="item.logo" alt="" class="w-60px h-60px">
      </div>
      <p>{{ item.text }}</p>
      <div v-if="item.versions && selectedRef.version" class="mt-4px">
        <el-select v-model="selectedRef.version" size="small">
          <el-option v-for="v in item.versions" :key="v.value" :label="v.label" :value="v.value" />
        </el-select>
      </div>
    </div>
  </div>

  <div class="mt-20px" :class="{ 'text-right': exportType === 'collection' }">
    <el-button type="primary" @click="handleExport(selectedRef)">
      {{ $t('app.common.export') }}
    </el-button>
  </div>
</template>
