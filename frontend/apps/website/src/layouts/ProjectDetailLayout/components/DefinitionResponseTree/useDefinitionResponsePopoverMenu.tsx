import AcTree from '@/components/AcTree'
import Node from '@/components/AcTree/model/node'
import { CollectionNode } from '@/typings/project'
import { ActiveNodeInfo } from '@/typings/common'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import NProgress from 'nprogress'
import { Menu } from '@/components/typings'
import useDefinitionResponseStore from '@/store/definitionResponse'
import { useActiveTree } from './useActiveTree'
import { useParams } from '@/hooks/useParams'
import { createDefaultResponseDefinition } from '@/views/document/components/createDefaultDefinition'
import { useGoPage } from '@/hooks/useGoPage'
import { useI18n } from 'vue-i18n'
import { ElCheckbox } from 'element-plus'
import { hasRefInSchema } from '@/commons'
/**
 * 目录弹层菜单逻辑
 * @param treeIns 目录树
 */
export const useDefinitionResponsePopoverMenu = (treeIns: Ref<InstanceType<typeof AcTree>>) => {
  const { t } = useI18n()

  const definitionResponseStore = useDefinitionResponseStore()
  const { project_id } = useParams()
  const { activeNode, reactiveNode } = useActiveTree(treeIns)
  const { goResponseEditPage } = useGoPage()

  const SCHEMA_MENUS: Menu[] = [{ text: t('app.common.delete'), onClick: () => onDeleteMenuClick() }]
  const popoverMenus = ref<Array<Menu>>(SCHEMA_MENUS)
  const popoverRefEl = ref<Nullable<HTMLElement>>(null)
  const isShowPopoverMenu = ref(false)
  const activeNodeInfo = ref<Nullable<ActiveNodeInfo>>({ node: undefined, id: undefined })

  const onPopoverRefIconClick = (e: Event, node?: Node) => {
    popoverMenus.value = SCHEMA_MENUS
    popoverRefEl.value = e.currentTarget as HTMLElement
    activeNodeInfo.value = { node, id: node?.data?.id }
    isShowPopoverMenu.value = true
  }

  /**
   * 删除模型
   */
  const onDeleteMenuClick = async () => {
    const node = unref(activeNodeInfo)?.node as Node
    const data = node?.data as any
    const tree = unref(treeIns)

    const isUnref = ref(1)

    AsyncMsgBox({
      title: t('app.common.deleteTip'),
      content: () => (
        <div>
          <div class="break-all mb-4px">{t('app.definitionResponse.popoverMenus.confirmDeleteDefinitionResponse', [data.name])}</div>
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
            对引用此相应的内容解引用
          </ElCheckbox>
        </div>
      ),
      onOk: async () => {
        NProgress.start()
        try {
          await definitionResponseStore.deleteDefinition(project_id as string, data.id, isUnref.value)
          tree.remove(node)
          activeNodeInfo.value = null
          reactiveNode()
          // directoryTree.reactiveNode && directoryTree.reactiveNode()
        } catch (error) {
        } finally {
          NProgress.done()
        }
      },
    })
  }

  /**
   * 复制
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
   * 创建
   */
  const onCreateMenuClick = async () => {
    const node = unref(activeNodeInfo)?.node as Node
    const tree = unref(treeIns)
    const newDefinition: any = createDefaultResponseDefinition({ name: t('app.definitionResponse.popoverMenus.unnamedDefinitionResponse') })

    try {
      NProgress.start()
      const newNode: any = await definitionResponseStore.createDefinition({ project_id, ...newDefinition })
      await nextTick()
      tree.setCurrentKey(newNode.id)
      goResponseEditPage(newNode.id)
      activeNode(newNode.id)
    } finally {
      NProgress.done()
    }
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
    onCreateMenuClick,
  }
}
