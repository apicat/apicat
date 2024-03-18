import NProgress from 'nprogress'
import { useI18n } from 'vue-i18n'
import type { ToggleHeading } from '@apicat/components'
import { ElCheckbox } from 'element-plus'
import type AcTreeWrapper from '../AcTreeWrapper'
import { usePopover } from '../usePopover'
import { PopoverMoreMenuType } from '../../constants'
import { useSelectedNode } from './useSelectedNode'
import type { Menu } from '@/components/typings'
import type Node from '@/components/AcTree/model/node'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { SchemaTypeEnum } from '@/commons'
import { useParams } from '@/hooks/useParams'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import { useProjectLayoutContext } from '@/layouts/ProjectDetailLayout/composables/useProjectLayoutContext'
import { injectPagesMode } from '@/layouts/ProjectDetailLayout/composables/usePagesMode'
import type { PageModeCtx } from '@/views/composables/usePageMode'

export function useMenus(
  treeIns: Ref<InstanceType<typeof AcTreeWrapper> | undefined>,
  toggleHeadingRef: Ref<InstanceType<typeof ToggleHeading> | undefined>,
) {
  const { t } = useI18n()
  const { projectID } = useParams()
  const { popoverRefEl, isShowPopoverMenu, show } = usePopover()
  const activeNodeInfo = ref<Node>()
  const schemaStore = useDefinitionSchemaStore()
  const { selectedNodeWithGoPage, reselectNode } = useSelectedNode(treeIns, toggleHeadingRef)
  const projectContext = useProjectLayoutContext()
  const { switchToWriteMode } = injectPagesMode('schema') as PageModeCtx

  // click wrapper function to check node is required
  function clickWrapper(fn: (node: Node) => void) {
    return () => {
      if (!activeNodeInfo.value)
        return
      fn(activeNodeInfo.value)
    }
  }

  /**
   * delete schema
   */
  function handleDelete(node: Node) {
    const data = node.data! as Definition.SchemaNode
    const isDir = data.type === SchemaTypeEnum.Category
    const isUnref = ref(true)

    AsyncMsgBox({
      title: t('app.common.deleteTip'),
      confirmButtonText: t('app.common.delete'),
      content: () => (
        <div>
          <div class="break-all">
            {isDir
              ? t('app.interface.popoverMenus.confirmDeleteGroup', [data.name])
              : t('app.interface.popoverMenus.confirmDeleteInterface', [data.name])}
          </div>
          {!isDir && (
            <ElCheckbox
              size="small"
              style={{ fontWeight: 'normal' }}
              modelValue={isUnref.value}
              onUpdate:modelValue={(val: any) => {
                isUnref.value = val
              }}>
              {t('app.schema.tips.unref')}
            </ElCheckbox>
          )}
        </div>
      ),
      onOk: async () => {
        try {
          await schemaStore.deleteSchema(projectID, data as Definition.SchemaNode, isDir ? false : isUnref.value)
          treeIns.value?.remove(node)
          reselectNode()
        }
        catch (error) {
          //
        }
      },
    })
  }

  /**
   * create schema
   */
  async function handleCreate(node?: Node) {
    try {
      NProgress.start()
      const data = await schemaStore.createSchema(projectID!, {
        name: t('app.schema.unnamedSchema'),
        type: SchemaTypeEnum.Schema,
        parentID: node?.key,
      })

      toggleHeadingRef.value?.expand()
      treeIns.value?.append(data, node!)
      selectedNodeWithGoPage(data)
      switchToWriteMode()
    }
    finally {
      NProgress.done()
    }
  }

  /**
   * create category
   */
  async function handleCreateCategroy(node?: Node) {
    try {
      NProgress.start()
      const data = await schemaStore.createSchema(projectID!, {
        name: t('app.schema.unnamedCategory'),
        type: SchemaTypeEnum.Category,
        parentID: node?.key,
      })

      toggleHeadingRef.value?.expand()
      const source: Definition.SchemaNode = node?.data || ({ items: schemaStore.schemas } as any)

      if (!source.items || !source.items.length)
        treeIns.value?.append(data, node!)
      else if (source.items && source.items.length)
        treeIns.value?.insertBefore(data, source.items[0] as any)

      await nextTick()
      treeIns.value?.rename(data)
    }
    finally {
      NProgress.done()
    }
  }

  /**
   * copy schema
   */
  async function handleCopy(node: Node) {
    try {
      NProgress.start()
      const data = await schemaStore.copySchema(projectID!, node.data as Definition.SchemaNode)
      treeIns.value?.insertAfter(data, node.key)
      selectedNodeWithGoPage(data)
      switchToWriteMode()
    }
    finally {
      NProgress.done()
    }
  }

  /**
   * rename menu click
   */
  function onRenameMenuClick(node: Node) {
    treeIns.value?.rename(node)
  }

  /**
   * rename schema
   */
  async function handleRename(node: Node, newName: string, oldName: string) {
    const data = node.data! as Definition.SchemaNode
    data.name = newName

    try {
      NProgress.start()
      await schemaStore.renameSchemaCategory(projectID!, { id: data.id, name: newName } as Definition.SchemaNode)
    }
    catch (error) {
      // 失败还原
      data.name = oldName
    }
    finally {
      NProgress.done()
    }
  }

  /**
   * AI
   */
  function onAICreate(node?: Node) {
    projectContext.handleAICreateSchema(async ({ prompt, showLoading, hideLoading }) => {
      try {
        showLoading()
        const data = await schemaStore.createSchemaWithAI(projectID!, {
          parentID: node?.key,
          prompt,
        })

        toggleHeadingRef.value?.expand()
        treeIns.value?.append(data, node!)
        selectedNodeWithGoPage(data)
        switchToWriteMode()

        projectContext.aiPromptDialogRef.value?.hide()
      }
      catch (error) {
        //
      }
      finally {
        hideLoading()
      }
    })
  }

  // 根目录添加菜单
  const ROOT_MENUS: Menu[] = [
    {
      text: t('app.schema.popoverMenus.newSchema'),
      icon: 'ac-model',
      onClick: () => handleCreate(activeNodeInfo.value),
    },
    {
      text: t('app.schema.popoverMenus.aiGenerateSchema'),
      icon: 'ac-zhinengyouhua',
      onClick: () => onAICreate(activeNodeInfo.value),
    },
    {
      text: t('app.schema.popoverMenus.newGroup'),
      icon: 'ac-folder',
      onClick: () => handleCreateCategroy(activeNodeInfo.value),
    },
  ]

  // 目录添加菜单
  const DIR_ADD_OPERATE_MENUS: Menu[] = ROOT_MENUS.concat([])

  // 目录更多菜单
  const DIR_MORE_OPERATE_MENUS_WITH_NODE: Menu[] = [
    { text: t('app.common.rename'), onClick: clickWrapper(onRenameMenuClick), icon: 'ac-edit' },
  ]
  const DIR_MORE_OPERATE_MENUS_EMPTY: Menu[] = DIR_MORE_OPERATE_MENUS_WITH_NODE.concat([
    { text: t('app.common.delete'), onClick: clickWrapper(handleDelete), icon: 'ac-trash' },
  ])

  // 文档更多菜单
  const DOC_MORE_OPERATE_MENUS: Menu[] = [
    { text: t('app.common.copy'), onClick: clickWrapper(handleCopy), icon: 'ac-copy' },
    { text: t('app.common.delete'), onClick: clickWrapper(handleDelete), icon: 'ac-trash' },
  ]

  // const popoverMenus = ref<Array<Menu>>([])
  const moreMenuTypeRef = ref<PopoverMoreMenuType>()
  const popoverMenus = computed<Array<Menu>>(() => {
    const node = activeNodeInfo.value
    const moreMenuType = moreMenuTypeRef.value
    // root 目录
    if (!node)
      return ROOT_MENUS

    // 目录菜单
    if (node && node.data?.type === SchemaTypeEnum.Category) {
      // 目录 添加
      if (moreMenuType === PopoverMoreMenuType.ADD) {
        const copyMenus = DIR_ADD_OPERATE_MENUS.concat([])
        // 层级限制
        if (node!.level >= 5)
          copyMenus.pop()

        return copyMenus
      }

      // 目录 更多
      if (moreMenuType === PopoverMoreMenuType.MORE) {
        const a = node.getChildren()
        return a && a.length > 0 ? DIR_MORE_OPERATE_MENUS_WITH_NODE : DIR_MORE_OPERATE_MENUS_EMPTY
      }
    }
    // 文档操作
    else {
      return DOC_MORE_OPERATE_MENUS
    }
    return []
  })

  // 节点菜单点击
  const onPopoverIconClick = (e: Event, node?: Node, moreMenuType?: PopoverMoreMenuType) => {
    show(e.currentTarget as HTMLElement)
    activeNodeInfo.value = node
    moreMenuTypeRef.value = moreMenuType
  }

  return {
    toggleHeadingRef,
    popoverMenus,
    popoverRefEl,
    isShowPopoverMenu,
    onPopoverIconClick,
    handleRename,
  }
}
