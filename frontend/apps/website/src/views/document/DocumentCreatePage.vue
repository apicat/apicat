<template>
  <div class="ac-header-operate">
    <div class="ac-header-operate__main">
      <p>
        <i class="ac-iconfont"></i>
        {{ isSaving ? '保存中...' : '已保存在云端' }}
      </p>
    </div>

    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="onPreviewBtnClick">预览</el-button>
    </div>
  </div>

  <HttpDocumentEditor v-model="httpDoc" />
</template>

<script setup lang="ts">
import { HttpDocument } from '@/typings'
import { createHttpDocument } from '@/views/document/components/createHttpDocument'
import { createCollection, updateCollection } from '@/api/collection'
import HttpDocumentEditor from './components/HttpDocumentEditor.vue'
import { useParams } from '@/hooks/useParams'
import { debounce, isEmpty } from 'lodash-es'
import { ElMessage } from 'element-plus'
import { useGoPage } from '@/hooks/useGoPage'
const { project_id } = useParams()
const isSaving = ref(false)
const httpDoc: Ref<HttpDocument> = ref(createHttpDocument())
const directoryTree: any = inject('directoryTree')
const { goDocumentDetailPage } = useGoPage()
let newDocId: any = null

watch(
  httpDoc,
  debounce(async () => {
    const data: any = { ...unref(httpDoc) }
    if (isEmpty(data.title)) {
      ElMessage.error('请输入文档标题')
      return
    }
    isSaving.value = true
    try {
      data.content = JSON.stringify(data.content)
      // update
      if (httpDoc.value.id) {
        const { id: collection_id, ...rest } = data
        await updateCollection({ project_id, collection_id, ...rest })
      } else {
        const res: any = await createCollection({ project_id, ...data })
        httpDoc.value.id = res.id
        newDocId = res.id
        directoryTree.createNodeByData(res)
      }
    } finally {
      isSaving.value = false
    }
  }, 400),
  {
    deep: true,
  }
)

const onPreviewBtnClick = () => {
  if (!newDocId) {
    return
  }

  goDocumentDetailPage(newDocId)
}
</script>
