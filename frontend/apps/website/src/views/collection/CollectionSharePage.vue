<script setup lang="ts">
import '@apicat/editor/dist/style.css'
import { useNamespace } from '@apicat/hooks'
import DocumentVerification from '../share/DocumentVerification.vue'
import { useCollectionShare } from './useCollectionShare'

const ns = useNamespace('document')
const AcEditor = defineAsyncComponent(() => import('@apicat/editor'))

const {
  schemas,
  responses,
  parameters,
  acEditorOptions,
  collectionInfo,
  collectionStatus,
  publicID,
  hideVerification,
  loading,
  onVerifyCodeInputSuccess,
} = useCollectionShare()
</script>

<template>
  <div v-loading="loading">
    <div v-if="hideVerification">
      <div class="flex w-full ac-top-nav shadow-black h-13 pl-30px pr-30px bg-white">
        <AcLogo class="ac-top-nav-logo" href="/" :size="90" />
      </div>
      <main class="flex flex-col min-h-screen bg-gray-100 p-30px pt90px">
        <div class="flex-1 bg-white p-30px">
          <h1>
            {{ collectionInfo?.title }}
          </h1>
          <AcEditor
            v-if="collectionInfo"
            :class="[ns.b()]"
            readonly
            :content="collectionInfo.content!"
            :schemas="schemas"
            :responses="responses"
            :parameters="parameters"
            :options="acEditorOptions" />
        </div>
      </main>
      <ac-backtop :bottom="100" :right="100" />
    </div>
    <DocumentVerification
      v-else
      v-model:visible="hideVerification"
      :project-i-d="collectionStatus?.projectID"
      :collection-i-d="collectionStatus?.collectionID"
      :public-i-d="publicID"
      @update:visible="onVerifyCodeInputSuccess" />
  </div>
</template>

<style scoped>
h1 {
  font-size: 28px;
  font-weight: 500;
  word-break: break-all;
  color: #27272a;
  outline: none;
  width: 100%;
  line-height: 45px;
  margin-bottom: 10px;
}
:deep(.ac-top-nav) {
  position: fixed;
  z-index: 100;
}
:deep(.ac-document) {
  padding: 0;
  min-height: unset;
}

.ac-top-nav {
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.2);
}
</style>
