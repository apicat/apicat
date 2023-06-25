import AcTree from '@/components/AcTree'
import createHttpDocIcon from '@/assets/images/doc-http@2x.png'
import Node from '@/components/AcTree/model/node'
import { Menu } from '@/components/typings'
import { CollectionNode } from '@/typings/project'
import { DocumentTypeEnum } from '@/commons/constant'
import { useRenameInput } from './useRenameInput'
import { ActiveNodeInfo } from '@/typings/common'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useDocumentStore, extendDocTreeFeild } from '@/store/document'
import { storeToRefs } from 'pinia'
import NProgress from 'nprogress'
import { useActiveTree } from './useActiveTree'
import { copyCollection, createCollection, deleteCollection } from '@/api/collection'
import { useParams } from '@/hooks/useParams'
import { createHttpDocument } from '@/views/document/components/createHttpDocument'
import { useGoPage } from '@/hooks/useGoPage'
import { useI18n } from 'vue-i18n'
import AIPromptModal from '../AIGenerateDocumentModal.vue'
import AcIconBIRobot from '~icons/bi/robot'

/**
 * hover 后更多菜单类型
 */
export const enum PopoverMoreMenuType {
  ADD,
  MORE,
}

let index = 1

/**
 * 目录弹层菜单逻辑
 * @param treeIns 目录树
 */
export const useDocumentPopoverMenu = (treeIns: Ref<InstanceType<typeof AcTree>>, aiPromptModalRef: Ref<InstanceType<typeof AIPromptModal>>) => {
  const { t } = useI18n()
  const popoverMenus = ref<Array<Menu>>([])
  const popoverMenuSize = ref('small')
  const popoverRefEl = ref<Nullable<HTMLElement>>(null)
  const isShowPopoverMenu = ref(false)
  const activeNodeInfo = ref<Nullable<ActiveNodeInfo>>({ node: undefined, id: undefined })
  const { onRenameMenuClick, ...otherInputLogic } = useRenameInput(activeNodeInfo)
  const documentStore = useDocumentStore()
  const { apiDocTree } = storeToRefs(documentStore)
  const { activeNode, reactiveNode } = useActiveTree(treeIns)
  const { project_id } = useParams()
  const { goDocumentEditPage } = useGoPage()
  const schemaTree = inject('schemaTree') as any
  const exportModal = inject('exportModal') as any

  const ROOT_MENUS: Menu[] = [
    { text: t('app.interface.popoverMenus.aiGenerateInterface'), elIcon: markRaw(AcIconBIRobot), onClick: () => onShowAIPromptModal() },
    { text: t('app.interface.popoverMenus.newInterface'), image: createHttpDocIcon, onClick: () => onCreateDocMenuClick() },
    { text: t('app.interface.popoverMenus.newGroup'), icon: 'ac-fenzu', onClick: () => onCreateDirMenuClick() },
  ]

  const DIR_ADD_OPERATE_MENUS: Menu[] = ROOT_MENUS.concat([])

  const DIR_MORE_OPERATE_MENUS: Menu[] = [
    {
      text: t('app.common.reanme'),
      onClick: () => onRenameMenuClick(),
    },
    {
      text: t('app.common.delete'),
      onClick: () => onDeleteMenuClick(),
    },
  ]

  const DOC_MORE_OPERATE_MENUS: Menu[] = [
    { text: t('app.common.copy'), onClick: () => onCopyMenuClick() },
    { text: t('app.common.export'), onClick: () => onExportMenuClick() },
    { text: t('app.common.delete'), onClick: () => onDeleteMenuClick() },
  ]

  const onPopoverRefIconClick = (e: Event, node?: Node, moreMenuType?: PopoverMoreMenuType) => {
    popoverRefEl.value = e.currentTarget as HTMLElement
    activeNodeInfo.value = { node, id: node?.data?.id }
    isShowPopoverMenu.value = true
    popoverMenuSize.value = 'small'

    // 顶级添加菜单
    if (!node) {
      popoverMenus.value = ROOT_MENUS
    }

    // 目录相关操作
    if (node && node?.data?.type === DocumentTypeEnum.DIR) {
      // 添加
      if (moreMenuType === PopoverMoreMenuType.ADD) {
        const copyMenus = DIR_ADD_OPERATE_MENUS.concat([])
        // 层级限制
        if (node!.level >= 5) {
          copyMenus.pop()
        }
        popoverMenus.value = copyMenus
      }

      // 更多
      if (moreMenuType === PopoverMoreMenuType.MORE) {
        popoverMenus.value = DIR_MORE_OPERATE_MENUS
      }
    }

    // 文档操作
    if (node && node?.data?.type !== DocumentTypeEnum.DIR) {
      popoverMenuSize.value = 'thin'
      popoverMenus.value = DOC_MORE_OPERATE_MENUS
    }
  }

  /**
   * 删除分类或文档
   */
  const onDeleteMenuClick = async () => {
    const tree = unref(treeIns)
    const node = unref(activeNodeInfo)?.node as Node
    const data = node?.data as CollectionNode
    const isDir = data.type === DocumentTypeEnum.DIR

    AsyncMsgBox({
      title: t('app.common.deleteTip'),
      content: (
        <div class="break-all">
          {isDir ? t('app.interface.popoverMenus.confirmDeleteGroup', [data.title]) : t('app.interface.popoverMenus.confirmDeleteInterface', [data.title])}
        </div>
      ),
      onOk: async () => {
        try {
          NProgress.start()
          await deleteCollection(project_id as string, data.id)
          tree.remove(node)
          reactiveNode()
          schemaTree.reactiveNode && schemaTree.reactiveNode()
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
      const newDoc: any = await copyCollection(project_id as string, data.id)
      tree.insertAfter(extendDocTreeFeild(newDoc), node)
    } finally {
      NProgress.done()
    }
  }

  /**
   * 创建分类
   */
  const onCreateDirMenuClick = async () => {
    const node = unref(activeNodeInfo)?.node as Node
    const source = node?.data as CollectionNode
    const tree = unref(treeIns)
    const data: any = { title: t('app.interface.popoverMenus.newGroup') + index++, type: DocumentTypeEnum.DIR }
    if (source && source.id) {
      data.parent_id = source.id
    }

    try {
      NProgress.start()
      const newNode: any = await createCollection({ project_id, ...data })
      const newData = extendDocTreeFeild(newNode, DocumentTypeEnum.DIR)
      if (!node) {
        apiDocTree.value.unshift(newData)
      } else {
        if (!source.items || !source.items.length) {
          tree.append(newData, node)
        } else {
          tree.insertBefore(newData, source.items[0])
        }
      }
      await nextTick()
      const parentNode = tree.getNode(source)
      parentNode && (parentNode.expanded = true)
      const current: Node = tree.getNode(newData.id)
      ;(current.data as CollectionNode)._extend!.isEditable = true
      setTimeout(() => otherInputLogic.renameInputFocus(), 100)
    } finally {
      NProgress.done()
    }
  }

  /**
   * 创建文档
   */
  const onCreateDocMenuClick = async () => {
    const node = unref(activeNodeInfo)?.node as Node
    const source = node?.data as CollectionNode
    const tree = unref(treeIns)
    const newDoc: any = createHttpDocument({ title: t('app.interface.popoverMenus.unnamedInterface') })
    newDoc.content = JSON.stringify(newDoc.content)
    const parent_id = !node ? 0 : source.id

    try {
      NProgress.start()
      const newNode: any = await createCollection({ project_id, parent_id, ...newDoc })
      const newData = extendDocTreeFeild(newNode)

      // root
      if (!node) {
        apiDocTree.value.unshift(newData)
      } else {
        if (!source.items || !source.items.length) {
          tree.append(newData, node)
        } else {
          tree.insertBefore(newData, source.items[0])
        }
      }

      await nextTick()
      tree.setCurrentKey(newNode.id)
      const parentNode = tree.getNode(source)
      parentNode && (parentNode.expanded = true)
      goDocumentEditPage(newNode.id)
      activeNode(newNode.id)
    } finally {
      NProgress.done()
    }
  }

  const createNodeByData = (data: any) => {
    const tree = unref(treeIns)
    const parentId = documentStore.tempCreateDocParentId
    const parentNode = parentId ? tree.getNode(parentId) : tree.root
    tree.append(extendDocTreeFeild(data), parentNode)
    nextTick(() => {
      tree.setCurrentKey(data.id)
      parentNode && (parentNode.expanded = true)
      activeNode(data.id)
    })
  }

  /**
   * 打开AI modal
   */
  const onShowAIPromptModal = () => {
    const node = unref(activeNodeInfo)?.node as Node
    const source = node?.data as CollectionNode
    const parent_id = !node ? 0 : source.id
    aiPromptModalRef.value.show({ parent_id })
  }

  /**
   * 导出
   */
  const onExportMenuClick = () => {
    const tree = unref(treeIns)
    const node = unref(activeNodeInfo)?.node as Node
    const data = node?.data as CollectionNode
    exportModal.exportDocument(project_id as string, data.id)
  }

  onClickOutside(popoverRefEl, () => {
    popoverRefEl.value = null
    isShowPopoverMenu.value = false
    activeNodeInfo.value!.id = undefined
  })

  onUnmounted(() => {
    index = 1
  })

  return {
    popoverRefEl,
    popoverMenus,
    isShowPopoverMenu,
    activeNodeInfo,
    popoverMenuSize,

    onPopoverRefIconClick,
    createNodeByData,
    ...otherInputLogic,
  }
}
