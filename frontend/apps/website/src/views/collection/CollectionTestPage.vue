<script setup lang="ts">
import { ElButton, ElDrawer, ElPopover, ElTable, ElTableColumn } from 'element-plus'
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n'
import type { PopoverInstance } from 'element-plus'
import { delay } from '@apicat/shared'
import type { Menu } from '@/components/typings'
import PromptTextarea from '@/components/PromptTextarea.vue'
import { apiCreateTestCase, apiGetTestCaseList, apiRegenTestCaseList } from '@/api/project/collection'

const props = defineProps<{ projectID: string; collectionID: string }>()
// 详情
const CollectionTestCaseDetail = defineAsyncComponent(() => import('./CollectionTestCaseDetail.vue'))
const collectionIDRef = toRef(props, 'collectionID')

const { t } = useI18n()
const list = ref<ProjectAPI.TestCase[]>([])

const prompt = ref('')

const visible = ref(false)
const generating = ref(false)
const tableLoading = ref(false)
const isShowPrompt = ref(false)
const isStartLoop = ref(false)

const promptTextareaRef = ref<InstanceType<typeof PromptTextarea>>()
const menusPopoverRef = ref<PopoverInstance>()
const generateBtnRef = ref()
const generateBtnTitle = computed(() => {
  if (generating.value)
    return t('app.project.collection.test.btn.loading')
  if (list.value.length)
    return t('app.project.collection.test.btn.continueGen')
  else return t('app.project.collection.test.btn.generate')
})

// 当前正在请求的AbortController
const currentAbortController = ref<AbortController>()
let loopController: ReturnType<typeof createLoop> | undefined

const menus: Array<Menu> = [
  {
    refText: generateBtnTitle,
    icon: 'ac-zhinengyouhua',
    onClick: onGeneratBtnClick,
  },
  {
    text: t('app.project.collection.test.btn.ai'),
    icon: 'ac-edit',
    onClick: onShowPromptBtnClick,
  },
  {
    text: t('app.project.collection.test.btn.regen'),
    icon: 'ac-reload',
    onClick: onRegenBtnClick,
  },
]

// 显示el-drawer
async function show() {
  reset()

  visible.value = true
  await refreshTestCases()
}

// 刷新测试用例列表
async function refreshTestCases() {
  tableLoading.value = true
  try {
    await getTestCaseList()
  }
  finally {
    tableLoading.value = false
  }
}

// 隐藏el-drawer
function hide() {
  visible.value = false
}

// el-drawer关闭事件
function onClosed() {
  reset()
}

// 重置状态
function reset() {
  isStartLoop.value = false
  currentAbortController.value?.abort()
  currentAbortController.value = undefined
  list.value = []
  prompt.value = ''
  isShowPrompt.value = false
  generating.value = false
  tableLoading.value = false
}

// 显示AI Prompt输入框
async function onShowPromptBtnClick() {
  isShowPrompt.value = true
}

// 点击重新生成所有测试用例按钮
async function onRegenBtnClick() {
  generating.value = true
  try {
    currentAbortController.value = new AbortController()
    await apiRegenTestCaseList(props.projectID, collectionIDRef.value, { signal: currentAbortController.value.signal })
    list.value = []
    isStartLoop.value = true
  }
  catch (error) {
    generating.value = false
  }
}

// 点击开始生成测试用例按钮
async function onGeneratBtnClick() {
  await generateTestcase()
}

// 获取测试用例列表
async function getTestCaseList() {
  currentAbortController.value = new AbortController()
  const result = await apiGetTestCaseList(props.projectID, collectionIDRef.value, {
    signal: currentAbortController.value.signal,
  })
  if (result) {
    generating.value = result.generating
    isStartLoop.value = result.generating
    list.value = result.records
  }
}

// 创建轮训
function createLoop(): { start: () => void; end: () => void } {
  let isLoop = false

  // 轮训测试用例列表
  async function run() {
    try {
      while (isLoop) {
        currentAbortController.value = new AbortController()
        const result = await apiGetTestCaseList(props.projectID, collectionIDRef.value, {
          signal: currentAbortController.value.signal,
        })
        if (result) {
          generating.value = result.generating
          isLoop = result.generating
          list.value = result.records
        }
        generating.value && (await delay(5000))
      }
    }
    catch (error) {
      end()
    }
  }

  function start() {
    isLoop = true
    generating.value = true
    run()
  }

  function end() {
    isLoop = false
    generating.value = false
  }

  return {
    start,
    end,
  }
}

// 隐藏AI输入框
function hidePromptTextarea() {
  isShowPrompt.value = false
  prompt.value = ''
}

// prompt输入框提交事件
async function onSubmitPrompt(promptStr: string) {
  if (!promptStr)
    return
  await generateTestcase(promptStr)
  hidePromptTextarea()
}

// 生成测试用例(直接生成&AI生成)
async function generateTestcase(prompt?: string) {
  generating.value = true
  try {
    currentAbortController.value = new AbortController()
    await apiCreateTestCase(props.projectID, collectionIDRef.value, prompt, false, {
      signal: currentAbortController.value.signal,
    })

    isStartLoop.value = true
  }
  catch (error) {
    generating.value = false
  }
}

// 删除测试用例事件
async function handleRemoveTestCase() {
  await refreshTestCases()
}

// 监测是否开始轮训
watch(isStartLoop, () => {
  if (isStartLoop.value) {
    loopController = createLoop()
    loopController.start()
  }
  else {
    loopController?.end()
    loopController = undefined
  }
})

// 进入测试用例详情logic
const detailRef = ref<InstanceType<typeof CollectionTestCaseDetail>>()
const inited = ref(false)
let resolveInit: any
async function intoTestCase(testCase: ProjectAPI.TestCase) {
  if (!inited.value) {
    const p = new Promise(res => (resolveInit = res))
    inited.value = true
    await p
    resolveInit = null
  }
  detailRef.value?.show(testCase.id)
}

// 浏览器历史记录监听
function popListener() {
  hide()
}

// 监听esc键,关闭AI输入框
function escListener(e: KeyboardEvent) {
  if (e.key === 'Escape' && isShowPrompt.value) {
    isShowPrompt.value = false
    prompt.value = ''
  }
}

function startListen() {
  document.addEventListener('keyup', escListener)
  window.addEventListener('popstate', popListener)
}
function stopListen() {
  document.removeEventListener('keyup', escListener)
  window.removeEventListener('popstate', popListener)
}

onMounted(() => startListen())
onUnmounted(() => stopListen())

defineExpose({
  show,
  hide,
})
</script>

<template>
  <ElDrawer
    v-model="visible"
    modal-class="TestCasedrawer"
    size="800px"
    destroy-on-close
    :show-close="false"
    :close-on-press-escape="!isShowPrompt"
    :with-header="false"
    @closed="onClosed"
  >
    <h4 class="mb-15px text-gray-title text-24px font-500">
      {{ $t('app.project.collection.test.title') }}
    </h4>

    <div v-loading="tableLoading" class="text-center roundBorder">
      <CollectionTestCaseDetail
        v-if="inited"
        ref="detailRef"
        :project-i-d="projectID"
        :collection-i-d="collectionID"
        @vue:before-mount="() => resolveInit && resolveInit()"
        @del="handleRemoveTestCase"
      />
      <ElTable :data="list" current-row-key="id" :row-style="{ cursor: 'pointer' }" @row-click="intoTestCase">
        <template #empty>
          <div class="leading-40px my-15px">
            <i class="ac-iconfont ac-empty text-45px" />
            <p>
              {{ $t('app.project.collection.test.table.empty') }}
            </p>
          </div>
        </template>
        <ElTableColumn type="index" width="50" align="center" />
        <ElTableColumn :label="$t('app.project.collection.test.table.name')">
          <template #default="{ row }">
            <p class="w-full truncate" :title="row.title">
              {{ row.title }}
            </p>
          </template>
        </ElTableColumn>
        <ElTableColumn
          prop="createdAt"
          :label="$t('app.project.collection.test.table.time')"
          width="200"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ dayjs(row.createdAt * 1000).format('LLL LT') }}
          </template>
        </ElTableColumn>
      </ElTable>

      <Transition :name="isShowPrompt ? 'slide-down' : 'slide-up'" mode="out-in">
        <ElButton
          v-if="!isShowPrompt"
          ref="generateBtnRef"
          class="mt-10px h-40px"
          :loading="generating"
          @click="onGeneratBtnClick"
        >
          {{ generateBtnTitle }}
        </ElButton>
        <PromptTextarea
          v-else
          ref="promptTextareaRef"
          v-model="prompt"
          class="mt-20px"
          prefix
          auto-focus
          :loading="generating"
          :placeholder="$t('app.project.collection.test.holder')"
          @submit="onSubmitPrompt"
        />
      </Transition>
    </div>
  </ElDrawer>

  <ElPopover
    v-if="list.length && !generating && !isShowPrompt"
    ref="menusPopoverRef"
    trigger="hover"
    virtual-triggering
    width="auto"
    transition="fade-fast"
    :virtual-ref="generateBtnRef"
    :show-arrow="false"
  >
    <PopperMenu size="small" class="clear-popover-space" :menus="menus" />
  </ElPopover>
</template>

<style lang="scss" scoped>
.roundBorder {
  :deep(.el-table__body) {
    border-collapse: separate;
    border-spacing: 0 5px;
  }
  :deep(.el-table--enable-row-hover .el-table__body tr td, .el-table__header thead tr th) {
    margin-top: 5px;
    border-bottom: none !important;
  }
  :deep(.el-table__inner-wrapper::before) {
    background-color: transparent !important;
  }
  :deep(.el-table--enable-row-hover .el-table__body tr td:first-child) {
    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
  }
  :deep(.el-table--enable-row-hover .el-table__body tr td:last-child) {
    border-bottom-right-radius: 5px;
    border-top-right-radius: 5px;
  }
  :deep(.el-table__header thead tr th:first-child) {
    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
  }
  :deep(.el-table__header thead tr th:last-child) {
    border-bottom-right-radius: 5px;
    border-top-right-radius: 5px;
  }
}
</style>

<style>
.TestCasedrawer {
  background: none;
}
</style>
