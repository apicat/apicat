import AcTree from '@/components/AcTree'
import Node from '@/components/AcTree/model/node'
import { CollectionNode } from '@/typings/project'
import { ActiveNodeInfo } from '@/typings/common'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import NProgress from 'nprogress'
import { Menu } from '@/components/typings'
import useDefinitionStore from '@/store/definition'
import { useActiveTree } from './useActiveTree'
import { deleteDefinition } from '@/api/definition'
import { useParams } from '@/hooks/useParams'
import createDefaultDefinition from '@/views/document/components/createDefaultDefinition'
import { useGoPage } from '@/hooks/useGoPage'
import AIGenerateSchemaModal from '../AIGenerateSchemaModal.vue'
import AIGenerateDocumentWithSchmeModal from '../AIGenerateDocumentWithSchmeModal.vue'
import AcIconBIRobot from '~icons/bi/robot'
import AcIconCarbonModelAlt from '~icons/carbon/model-alt'
/**
 * 目录弹层菜单逻辑
 * @param treeIns 目录树
 */
export const useSchemaPopoverMenu = (
  treeIns: Ref<InstanceType<typeof AcTree>>,
  aiPromptModalRef: Ref<InstanceType<typeof AIGenerateSchemaModal>>,
  aiGenerateDocumentWithSchemaModalRef: Ref<InstanceType<typeof AIGenerateDocumentWithSchmeModal>>
) => {
  const definitionStore = useDefinitionStore()
  const { project_id } = useParams()
  const { activeNode, reactiveNode } = useActiveTree(treeIns)
  const { goSchemaEditPage } = useGoPage()

  const ROOT_MENUS: Menu[] = [
    { text: 'AI生成模型', elIcon: markRaw(AcIconBIRobot), onClick: () => onShowAIPromptModal() },
    { text: '新建模型', elIcon: markRaw(AcIconCarbonModelAlt), onClick: () => onCreateSchemaMenuClick() },
  ]

  const SCHEMA_MENUS: Menu[] = [
    { text: 'AI生成接口', onClick: () => onCreateDocumentBySchema() },
    { text: '复制', onClick: () => onCopyMenuClick() },
    { text: '删除', onClick: () => onDeleteMenuClick() },
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
   * 删除分类或文档
   */
  const onDeleteMenuClick = async () => {
    const node = unref(activeNodeInfo)?.node as Node
    const data = node?.data as CollectionNode
    const tree = unref(treeIns)

    AsyncMsgBox({
      title: '删除提示',
      content: <div class="break-all">确定删除「{data.name}」模型吗？</div>,
      onOk: async () => {
        NProgress.start()
        try {
          await deleteDefinition(project_id as string, data.id)
          tree.remove(node)
          activeNodeInfo.value = null
          reactiveNode()
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
    const newDefinition: any = createDefaultDefinition({ name: 'Unnamed' })

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
