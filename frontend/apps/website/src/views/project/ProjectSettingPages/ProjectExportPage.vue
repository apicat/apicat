<template>
  <div class="flex gap-30px">
    <div class="flex flex-col items-center" v-for="item in exportList" :key="item.text" @click="selectedRef = item">
      <div class="border border-solid rounded cursor-pointer p-20px hover:border-blue-primary" :class="[{ 'border-blue-primary': selectedRef.type === item.type }]">
        <img :src="item.logo" alt="" class="w-60px h-60px" />
      </div>
      <p>{{ item.text }}</p>
    </div>
  </div>

  <el-button type="primary" class="mt-20px" @click="handelExport(selectedRef.type)">{{ $t('app.common.export') }}</el-button>
</template>
<script setup lang="ts">
import swaggerLogo from '@/assets/images/logo-swagger@2x.png'
import openApiLogo from '@/assets/images/logo-postman@2x.png'
import { ExportProjectTypes } from '@/commons/constant'
import { exportProject } from '@/api/project'
import { useProjectId } from '@/hooks/useProjectId'

const project_id = useProjectId()

const exportList = [
  { logo: swaggerLogo, text: 'Swagger 2.0', type: ExportProjectTypes.Swagger },
  { logo: openApiLogo, text: 'OpenAPI 3.0', type: ExportProjectTypes.OpenAPI },
]
const selectedRef = ref(exportList[0])

const handelExport = (type: ExportProjectTypes) => window.open(exportProject({ project_id, type }))
</script>
