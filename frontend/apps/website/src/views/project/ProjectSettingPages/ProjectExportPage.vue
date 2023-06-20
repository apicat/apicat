<template>
  <div class="flex gap-30px">
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

  <el-button type="primary" class="mt-20px" @click="handleExport(selectedRef)">{{ $t('app.common.export') }}</el-button>
</template>
<script setup lang="ts">
import swaggerLogo from '@/assets/images/logo-swagger@2x.png'
import openApiLogo from '@/assets/images/logo-openapis.svg'
import htmlLogo from '@/assets/images/logo-html@2x.png'
import mdLogo from '@/assets/images/logo-markdown@2x.png'
import { ExportProjectTypes } from '@/commons/constant'
import { exportProject } from '@/api/project'
import { useProjectId } from '@/hooks/useProjectId'

const project_id = useProjectId()

const exportList = [
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
  { logo: htmlLogo, text: 'HTML', type: ExportProjectTypes.HTML },
  { logo: mdLogo, text: 'Markdown', type: ExportProjectTypes.MARKDOWN },
]

const selectedRef: any = ref({
  type: exportList[0].type,
  version: 'openapi3.0.0',
})

const handleSelect = (selected: any, item: any) => {
  selected.type = item.type
}

const handleExport = (selected: any) => {
  let type = selected.type
  if (selected.type === ExportProjectTypes.OpenAPI) {
    type = selected.version
  }

  window.open(exportProject({ project_id, type }))
}
</script>
