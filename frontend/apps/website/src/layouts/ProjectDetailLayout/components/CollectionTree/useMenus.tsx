import NProgress from 'nprogress'
import { useI18n } from 'vue-i18n'
import type { ToggleHeading } from '@apicat/components'
import type AcTreeWrapper from '../AcTreeWrapper'
import { usePopover } from '../usePopover'
import { PopoverMoreMenuType } from '../../constants'
import { useSelectedNode } from './useSelectedNode'
import type { Menu } from '@/components/typings'
import type Node from '@/components/AcTree/model/node'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { CollectionTypeEnum } from '@/commons'
import { useCollectionsStore } from '@/store/collections'
import { useParams } from '@/hooks/useParams'
import { useProjectLayoutContext } from '@/layouts/ProjectDetailLayout/composables/useProjectLayoutContext'
import { injectPagesMode } from '@/layouts/ProjectDetailLayout/composables/usePagesMode'
import type { PageModeCtx } from '@/views/composables/usePageMode'

export function useMenus(
  treeIns: Ref<InstanceType<typeof AcTreeWrapper> | undefined>,
  toggleHeadingRef: Ref<InstanceType<typeof ToggleHeading> | undefined>,
) {
  const { t } = useI18n()
  const popoverMenus: Ref<Menu> = ref<Menu[]>([])
  const { projectID } = useParams()
  const { popoverRefEl, isShowPopoverMenu, show } = usePopover()
  const activeNodeInfo = ref<Node>()
  const collectionsStore = useCollectionsStore()
  const { selectedNodeWithGoPage, reselectNode } = useSelectedNode(treeIns, toggleHeadingRef)
  const projectContext = useProjectLayoutContext()
  const { switchToWriteMode } = injectPagesMode('collection') as PageModeCtx

  // click wrapper function to check node is required
  function clickWrapper(fn: (node: Node) => void) {
    return () => {
      if (!activeNodeInfo.value)
        return
      fn(activeNodeInfo.value)
    }
  }

  /**
   * 删除集合
   */
  function onDeleteCollection(node: Node) {
    const data = node.data!
    const isDir = data.type === CollectionTypeEnum.Dir

    AsyncMsgBox({
      title: t('app.common.deleteTip'),
      confirmButtonText: t('app.common.delete'),
      content: (
        <div class="break-all">
          {isDir
            ? t('app.interface.popoverMenus.confirmDeleteGroup', [data.title])
            : t('app.interface.popoverMenus.confirmDeleteInterface', [data.title])}
        </div>
      ),
      onOk: async () => {
        try {
          await collectionsStore.deleteCollection(projectID, data.id)
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
   * 创建集合
   */
  async function onCreateCollection(node?: Node) {
    try {
      NProgress.start()
      const data = await collectionsStore.createCollection(projectID!, {
        title: t('app.interface.unnamedInterface'),
        type: CollectionTypeEnum.Http,
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
   * 创建目录
   */
  async function onCreateDir(node?: Node) {
    try {
      NProgress.start()
      const data = await collectionsStore.createCollection(projectID!, {
        title: t('app.interface.unnamedCategory'),
        type: CollectionTypeEnum.Dir,
        parentID: node?.key,
      })

      toggleHeadingRef.value?.expand()

      const source: CollectionAPI.ResponseCollection = node?.data || ({ items: collectionsStore.collections } as any)

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
   * 复制集合
   */
  async function onCopyCollection(node: Node) {
    try {
      NProgress.start()
      const data = await collectionsStore.copyCollection(projectID!, node.data?.id)
      treeIns.value?.insertAfter(data, node.key)
      selectedNodeWithGoPage(data)
      switchToWriteMode()
    }
    catch (error) {
      //
    }
    finally {
      NProgress.done()
    }
  }

  /**
   * 点击重命名菜单
   */
  function onRenameMenuClick(node: Node) {
    treeIns.value?.rename(node)
  }

  /**
   * 重命名集合
   */
  async function handleRenameCollection(node: Node, newName: string, oldName: string) {
    const data = node.data!
    data.title = newName

    try {
      NProgress.start()
      await collectionsStore.renameCollection(projectID!, {
        id: data.id,
        title: newName,
      } as CollectionAPI.ResponseCollection)
    }
    catch (error) {
      // 失败还原
      data.title = oldName
    }
    finally {
      NProgress.done()
    }
  }
  /**
   * 导出集合
   */
  function onExportCollection(node?: Node) {
    projectContext.handleExportDocument!(projectID, node!.key)
  }

  /**
   * AI
   */
  function onAICreate(node?: Node) {
    projectContext.handleAICreateCollection(async ({ prompt, showLoading, hideLoading }) => {
      try {
        showLoading()
        const data = await collectionsStore.createCollectionWithAI(projectID!, {
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
      text: t('app.interface.popoverMenus.newInterface'),
      icon: 'ac-doc',
      onClick: () => onCreateCollection(activeNodeInfo.value),
    },
    {
      text: t('app.interface.popoverMenus.aiGenerateInterface'),
      icon: 'ac-zhinengyouhua',
      onClick: () => onAICreate(activeNodeInfo.value),
    },
    {
      text: t('app.interface.popoverMenus.newGroup'),
      icon: 'ac-folder',
      onClick: () => onCreateDir(activeNodeInfo.value),
    },
  ]

  // 目录添加菜单
  const DIR_ADD_OPERATE_MENUS: Menu[] = ROOT_MENUS.concat([])

  // 目录更多菜单
  const DIR_MORE_OPERATE_MENUS: Menu[] = [
    { text: t('app.common.rename'), onClick: clickWrapper(onRenameMenuClick), icon: 'ac-edit' },
    { text: t('app.common.delete'), onClick: clickWrapper(onDeleteCollection), icon: 'ac-trash' },
  ]

  // 文档更多菜单
  const DOC_MORE_OPERATE_MENUS: Menu[] = [
    { text: t('app.common.copy'), onClick: clickWrapper(onCopyCollection), icon: 'ac-copy' },
    { text: t('app.common.export'), onClick: clickWrapper(onExportCollection), icon: 'ac-import' },
    { text: t('app.common.delete'), onClick: clickWrapper(onDeleteCollection), icon: 'ac-trash' },
  ]

  // 节点菜单点击
  const onPopoverIconClick = (e: Event, node?: Node, moreMenuType?: PopoverMoreMenuType) => {
    show(e.currentTarget as HTMLElement)
    activeNodeInfo.value = node

    // root 目录
    if (!node) {
      popoverMenus.value = ROOT_MENUS
      return
    }

    // 目录菜单
    if (node && node.data?.type === CollectionTypeEnum.Dir) {
      // 目录 添加
      if (moreMenuType === PopoverMoreMenuType.ADD) {
        const copyMenus = DIR_ADD_OPERATE_MENUS.concat([])
        // 层级限制
        if (node!.level >= 5)
          copyMenus.pop()

        popoverMenus.value = copyMenus
      }

      // 目录 更多
      if (moreMenuType === PopoverMoreMenuType.MORE)
        popoverMenus.value = DIR_MORE_OPERATE_MENUS
    }
    // 文档操作
    else {
      popoverMenus.value = DOC_MORE_OPERATE_MENUS
    }
  }

  return {
    toggleHeadingRef,
    popoverMenus,
    popoverRefEl,
    isShowPopoverMenu,
    onPopoverIconClick,
    handleRenameCollection,
  }
}
