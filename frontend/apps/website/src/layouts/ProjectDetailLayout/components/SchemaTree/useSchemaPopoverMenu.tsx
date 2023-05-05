import AcTree from '@/components/AcTree'
import Node from '@/components/AcTree/model/node'
import { CollectionNode } from '@/typings/project'
import { ActiveNodeInfo } from '@/typings/common'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import NProgress from 'nprogress'
import { Menu } from '@/components/typings'
import useDefinitionStore from '@/store/definition'
import { useActiveTree } from './useActiveTree'
import { useParams } from '@/hooks/useParams'
import createDefaultDefinition from '@/views/document/components/createDefaultDefinition'
import { useGoPage } from '@/hooks/useGoPage'
import AIGenerateSchemaModal from '../AIGenerateSchemaModal.vue'
import AIGenerateDocumentWithSchmeModal from '../AIGenerateDocumentWithSchmeModal.vue'
import AcIconBIRobot from '~icons/bi/robot'
import AcIconCarbonModelAlt from '~icons/carbon/model-alt'
import { useI18n } from 'vue-i18n'
import { ElCheckbox, ElSwitch } from 'element-plus'
import { h } from 'vue'
/**
 * 目录弹层菜单逻辑
 * @param treeIns 目录树
 */
export const useSchemaPopoverMenu = (
  treeIns: Ref<InstanceType<typeof AcTree>>,
  aiPromptModalRef: Ref<InstanceType<typeof AIGenerateSchemaModal>>,
  aiGenerateDocumentWithSchemaModalRef: Ref<InstanceType<typeof AIGenerateDocumentWithSchmeModal>>
) => {
  const { t } = useI18n()

  const definitionStore = useDefinitionStore()
  const { project_id } = useParams()
  const { activeNode, reactiveNode } = useActiveTree(treeIns)
  const { goSchemaEditPage } = useGoPage()
  const directoryTree = inject('directoryTree') as any

  const ROOT_MENUS: Menu[] = [
    { text: t('app.schema.popoverMenus.aiGenerateSchema'), elIcon: markRaw(AcIconBIRobot), onClick: () => onShowAIPromptModal() },
    { text: t('app.schema.popoverMenus.newSchema'), elIcon: markRaw(AcIconCarbonModelAlt), onClick: () => onCreateSchemaMenuClick() },
  ]

  const SCHEMA_MENUS: Menu[] = [
    { text: t('app.interface.popoverMenus.aiGenerateInterface'), onClick: () => onCreateDocumentBySchema() },
    { text: t('app.common.copy'), onClick: () => onCopyMenuClick() },
    { text: t('app.common.delete'), onClick: () => onDeleteMenuClick() },
  ]
  const popoverMenus = ref<Array<Menu>>(SCHEMA_MENUS)
  const popoverRefEl = ref<Nullable<HTMLElement>>(null)
  const isShowPopoverMenu = ref(false)
  const activeNodeInfo = ref<Nullable<ActiveNodeInfo>>({ node: undefined, id: undefined })

  const onPopoverRefIconClick = (e: Event, node?: Node) => {
    popoverMenus.value = SCHEMA_MENUS

    // 顶级添加菜单
    if (!node) {
      popoverMenus.value = ROOT_MENUS
    }

    popoverRefEl.value = e.currentTarget as HTMLElement
    activeNodeInfo.value = { node, id: node?.data?.id }
    isShowPopoverMenu.value = true
  }

  /**
   * 删除模型
   */
  const onDeleteMenuClick = async () => {
    const node = unref(activeNodeInfo)?.node as Node
    const data = node?.data as CollectionNode
    const tree = unref(treeIns)
    const isUnref = ref(1)

    AsyncMsgBox({
      title: t('app.common.deleteTip'),
      content: () => (
        <div>
          <div class="break-all mb-4px">{t('app.interface.popoverMenus.confirmDeleteInterface', [data.name])}</div>
          <ElCheckbox
            size="small"
            style={{ fontWeight: 'normal' }}
            modelValue={isUnref.value}
            onUpdate:modelValue={(val: any) => {
              isUnref.value = val
            }}
            trueLabel={1}
            falseLabel={0}
          >
            对引用此模型的内容解引用
          </ElCheckbox>
        </div>
      ),
      onOk: async () => {
        NProgress.start()
        try {
          await definitionStore.deleteDefinition(project_id as string, data.id, isUnref.value)
          tree.remove(node)
          activeNodeInfo.value = null
          reactiveNode()
          directoryTree.reactiveNode && directoryTree.reactiveNode()
        } catch (error) {
        } finally {
          NProgress.done()
        }
      },
    })
  }

  /**
   * 复制文档
   */
  const onCopyMenuClick = async () => {
    const tree = unref(treeIns)
    const node = unref(activeNodeInfo)?.node as Node
    const data = node?.data as CollectionNode
    try {
      NProgress.start()
      await definitionStore.copyDefinition(project_id as string, data.id)
    } finally {
      NProgress.done()
    }
  }

  /**
   * 创建模型
   */
  const onCreateSchemaMenuClick = async () => {
    const node = unref(activeNodeInfo)?.node as Node
    const tree = unref(treeIns)
    const newDefinition: any = createDefaultDefinition({ name: t('app.schema.popoverMenus.unnamedSchema') })

    try {
      NProgress.start()
      const newNode: any = await definitionStore.createDefinition({ project_id, ...newDefinition })
      await nextTick()
      tree.setCurrentKey(newNode.id)
      goSchemaEditPage(newNode.id)
      activeNode(newNode.id)
    } finally {
      NProgress.done()
    }
  }

  /**
   * 打开AI modal
   */
  const onShowAIPromptModal = () => {
    aiPromptModalRef.value.show()
  }

  const onCreateDocumentBySchema = () => {
    const node = unref(activeNodeInfo)?.node as Node
    const data = node?.data as CollectionNode
    aiGenerateDocumentWithSchemaModalRef.value.show(data)
  }

  onClickOutside(popoverRefEl, () => {
    popoverRefEl.value = null
    isShowPopoverMenu.value = false
    activeNodeInfo.value!.id = undefined
  })

  return {
    popoverRefEl,
    popoverMenus,
    isShowPopoverMenu,
    activeNodeInfo,

    onPopoverRefIconClick,
    onCreateSchemaMenuClick,
  }
}
