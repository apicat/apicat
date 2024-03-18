<script setup lang="ts">
import '@/styles/github-markdown.css'
import 'highlight.js/styles/vs.min.css'
import { ElDivider, ElDrawer } from 'element-plus'
import { Marked } from 'marked'
import { markedHighlight } from 'marked-highlight'
import hljs from 'highlight.js'
import { useI18n } from 'vue-i18n'
import { apiDeleteTestCase, apiGetTestCaseDetail, apiReGenTestCase } from '@/api/project/collection'
import useApi from '@/hooks/useApi'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'

const props = defineProps<{
  projectID: string
  collectionID: string | number
}>()

const emit = defineEmits(['del', 'closed'])

const marked = new Marked(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext'
      return hljs.highlight(code, { language }).value
    },
  }),
)

const { t } = useI18n()
const testcaseID = ref(-1)
const visible = ref(false)
const data = ref<ProjectAPI.TestCaseDetail>()
const [loading, getTestcaseInfo] = useApi(apiGetTestCaseDetail)
async function getData() {
  try {
    const res = await getTestcaseInfo(props.projectID, props.collectionID, testcaseID.value)
    if (res)
      res.content = await marked.parse(res.content)

    data.value = res
  }
  catch (e) {
    // hide()
  }
}

const btnState = ref(true)
const [reLoading, regen] = useApi(apiReGenTestCase)
const prompt = ref('')
async function onSubmit() {
  loading.value = true
  try {
    const res = await regen(props.projectID, props.collectionID as number, testcaseID.value, prompt.value)
    prompt.value = ''
    if (res) {
      res.content = await marked.parse(res.content)
      data.value = res
    }
    btnState.value = true
  }
  finally {
    loading.value = false
  }
}
const [_, del] = useApi(apiDeleteTestCase)
const topRef = ref<HTMLElement>()

async function onDelete() {
  AsyncMsgBox({
    title: t('app.project.collection.test.detail.del.title'),
    message: t('app.project.collection.test.detail.del.tip'),
    confirmButtonClass: 'red',
    confirmButtonText: t('app.project.collection.test.detail.del.confirm'),
    cancelButtonText: t('app.common.cancel'),
    appendTo: topRef.value,
    customClass: 'hide-modal',
    onOk: async () => {
      await del(props.projectID, props.collectionID as number, testcaseID.value)
      emit('del', testcaseID.value)
      hide()
    },
  })
}

function reGen() {
  btnState.value = false
}

// 控制esc关闭和退出（配合btnState）
const promptRef = ref<HTMLElement>()
function escListener(e: KeyboardEvent) {
  const target = e.target as HTMLElement
  function cleanKey(e: KeyboardEvent) {
    if (e.shiftKey || e.ctrlKey || e.altKey || e.metaKey)
      return false
    else return true
  }
  if (e.key === 'Escape' && !btnState.value) {
    btnState.value = true
    prompt.value = ''
  }
  else if (cleanKey(e) && target.tagName !== 'INPUT' && target.tagName !== 'TEXTAREA') {
    if (e.key === 'd')
      onDelete()
    if (e.key === 'r')
      reGen()
  }
}
const addListener = () => document.addEventListener('keyup', escListener)
const removeListener = () => document.removeEventListener('keyup', escListener)

function show(id: number) {
  testcaseID.value = id
  getData()
  addListener()
  visible.value = true
}

function hide() {
  removeListener()
  visible.value = false
  testcaseID.value = -1
  data.value = undefined
}

// 阴影
const showShadow = ref(false)
const checkShadow = (function () {
  let timer: any
  return () => {
    if (timer)
      clearTimeout(timer)
    timer = setTimeout(() => {
      showShadow.value = topRef.value ? topRef.value.getBoundingClientRect().y < -10 : false
    }, 0)
  }
})()
function drawerOpened() {
  const scrollBody = document.querySelector('.TestCaseDetailDrawer .el-drawer__body')
  if (scrollBody) {
    scrollBody.addEventListener('scroll', checkShadow)
    checkShadow()
  }
}
function drawerClose() {
  const scrollBody = document.querySelector('.TestCaseDetailDrawer .el-drawer__body')
  if (scrollBody)
    scrollBody.removeEventListener('scroll', checkShadow)
}

function closed() {
  btnState.value = true
  emit('closed')
}

onUnmounted(() => removeListener())

defineExpose({
  show,
  hide,
})
</script>

<template>
  <ElDrawer
    v-model="visible"
    size="800px"
    modal-class="TestCaseDetailDrawer"
    destroy-on-close
    :with-header="false"
    :show-close="false"
    :close-on-press-escape="btnState"
    :append-to-body="true"
    @opened="drawerOpened"
    @close="drawerClose"
    @closed="closed"
  >
    <transition name="fade">
      <div v-if="loading" class="mask">
        <div class="el-loading-spinner">
          <svg class="circular" viewBox="0 0 50 50"><circle class="path" cx="25" cy="25" r="20" fill="none" /></svg>
        </div>
      </div>
    </transition>
    <div
      class="absolute w-full bg-white cursor-pointer top-0px pt-30px pb-20px z-100"
      :class="[{ shadow: showShadow }]"
      @click="hide"
    >
      <i class="leading-none ac-iconfont ac-back" style="padding: 20px; font-size: 16px" />
    </div>
    <div class="verflow-hidden" @scroll="checkShadow">
      <div ref="topRef" class="pt-80px pl-20px pr-20px pb-20px testcase-detail-content overscroll-y-auto min-h-100vh">
        <div v-if="!data" class="text-center leading-40px my-15px">
          <i class="ac-iconfont ac-empty text-45px" />
          <p>
            {{ $t('app.project.collection.test.table.empty') }}
          </p>
        </div>

        <div v-else>
          <h4 class="text-gray-title text-24px font-500">
            {{ data?.title }}
          </h4>
          <div style="position: relative">
            <div class="break-words mt-30px">
              <div class="markdown-body" v-html="data?.content" />
            </div>
          </div>

          <Transition :name="btnState ? 'slide-down' : 'slide-up'" mode="out-in">
            <div v-if="btnState" class="row">
              <div class="left">
                <ElButton link @click="reGen">
                  <i class="leading-none ac-iconfont ac-loop-inner" style="font-size: 16px" />
                  <span class="ml-5px"> {{ $t('app.project.collection.test.detail.regenerate') }} </span>
                </ElButton>
              </div>
              <div class="right">
                <ElButton link @click="onDelete">
                  <i class="leading-none ac-iconfont ac-trash" style="font-size: 16px" />
                </ElButton>
              </div>
            </div>
            <template v-else>
              <div>
                <ElDivider />
                <PromptTextarea
                  ref="promptRef"
                  v-model="prompt"
                  class="mt-20px"
                  prefix
                  empty-trigger
                  :placeholder="$t('app.project.collection.test.detail.holder')"
                  :loading="reLoading"
                  @vue:mounted="(e: any) => e.component.exposed.focus()"
                  @submit="onSubmit"
                />
              </div>
            </template>
          </Transition>
        </div>
      </div>
    </div>
  </ElDrawer>
</template>

<style lang="scss" scoped>
:deep(.el-divider) {
  border: 1px solid rgb(221, 221, 221);
  margin: 20px 0;
}

.row {
  margin-top: 1em;
  margin-bottom: 1em;
  display: flex;
  justify-content: space-between;
  width: 100%;
}
.left,
.right {
  display: flex;
  align-items: center;
  font-size: 30px;
}
.left {
  flex-grow: 1;
}
.right {
  justify-content: flex-start;
}
</style>

<style lang="scss">
.TestCaseDetailDrawer {
  background: none;
  .el-drawer {
    box-shadow: none;
  }
  .el-drawer__body {
    padding: 0;
  }
  .el-overlay {
    position: absolute;
    .el-overlay-message-box {
      position: absolute;
    }
  }

  .el-overlay {
    background-color: transparent;
  }
}
</style>
