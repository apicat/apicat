import type { InjectionKey } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Menu } from '@/components/typings'
import { Authority } from '@/commons/constant'
import General from '@/views/project/ProjectSettingPages/BaseInfoSetting.vue'
import Member from '@/views/project/ProjectSettingPages/ProjectMemberList.vue'
import Url from '@/views/project/ProjectSettingPages/ServerUrlSetting.vue'
import GlobalParam from '@/views/project/ProjectSettingPages/GlobalParametersSetting.vue'
import Export from '@/views/project/ProjectSettingPages/ProjectExportPage.vue'
import Trash from '@/views/project/ProjectSettingPages/ProjectTrashPage.vue'
import DeleteOrQuit from '@/views/project/ProjectSettingPages/DeleteOrQuit.vue'
import type Node from '@/components/AcTree/model/node'

export interface CurrentNodeContext {
  activeCollectionNode: (id: number) => void
  activeSchemaNode: (id: number) => void
  activeResponseNode: (id: number) => void
  currentActiveNode: Ref<{ id?: number; type: string }>
  activeCollectionKey: ComputedRef<number | undefined>
  activeSchemaKey: ComputedRef<number | undefined>
  activeResponseKey: ComputedRef<number | undefined>
}
export const CurrentNodeContextKey: InjectionKey<CurrentNodeContext> = Symbol('CurrentNodeContextKey')

export enum PopoverMoreMenuType {
  ADD,
  MORE,
}

/**
 * 项目布局中导航菜单
 */
export enum ProjectNavigateListEnum {
  General = 'general',
  Member = 'member',
  URL = 'url',
  GlobalParam = 'globalParam',
  ResponseParamsSetting = 'response',
  Export = 'export',
  ProjectShare = 'share',
  Trash = 'trash',
  Delete = 'delete',
  Quit = 'quit',
}

/**
 *
 * @returns use function in setup
 * { [key in ProjectNavigateListEnum]: { [key: string]: any }
 */
export function getProjectNavigateList(
  isModal: boolean,
  auth: Authority,
  overwrite?: any,
): {
    [key: string]: Menu
  } {
  const { t } = useI18n()
  let navs = {
    [ProjectNavigateListEnum.General]: {
      detailTitle: t('app.project.setting.basic.detailTitle'),
      text: t('app.project.setting.basic.title'),
      icon: 'ac-setting',
      sort: 100,
      component: General,
      auth: [Authority.Manage],
    },
    [ProjectNavigateListEnum.Member]: {
      detailTitle: t('app.project.setting.member.detailTitle'),
      text: t('app.project.setting.member.title'),
      icon: 'ac-members1',
      sort: 200,
      component: Member,
      auth: [Authority.Manage, Authority.Write, Authority.Read],
    },
    [ProjectNavigateListEnum.URL]: {
      detailTitle: t('app.project.setting.urls.detailTitle'),
      text: t('app.project.setting.urls.title'),
      icon: 'ac-url',
      sort: 300,
      component: Url,
      auth: [Authority.Manage, Authority.Write],
    },
    [ProjectNavigateListEnum.GlobalParam]: {
      detailTitle: t('app.project.setting.globalParam.detailTitle'),
      text: t('app.project.setting.globalParam.title'),
      icon: 'ac-canshuweihu',
      sort: 400,
      component: GlobalParam,
      auth: [Authority.Manage, Authority.Write],
    },
    [ProjectNavigateListEnum.ProjectShare]: {
      detailTitle: t('app.project.setting.share.detailTitle'),
      text: t('app.project.setting.share.title'),
      icon: 'ac-share',
      sort: 500,
      auth: [Authority.Manage, Authority.Write],
    },
    [ProjectNavigateListEnum.Export]: {
      detailTitle: t('app.project.setting.export.detailTitle'),
      text: t('app.project.setting.export.title'),
      icon: 'ac-download',
      sort: 600,
      component: Export,
      auth: [Authority.Manage, Authority.Write],
    },
    [ProjectNavigateListEnum.Trash]: {
      detailTitle: t('app.project.setting.trash.detailTitle'),
      text: t('app.project.setting.trash.title'),
      icon: 'ac-trash',
      sort: 700,
      component: Trash,
      auth: [Authority.Manage, Authority.Write],
    },
  } as any

  Object.keys(navs).forEach((key: any) => {
    navs[key].key = key
  })

  if (isModal) {
    navs = filterModel(navs, {
      [ProjectNavigateListEnum.Delete]: {
        detailTitle: t('app.project.setting.delete.btitle'),
        text: t('app.project.setting.delete.btitle'),
        icon: 'ac-delete',
        sort: 800,
        component: DeleteOrQuit,
        auth: [Authority.Manage],
      },
      [ProjectNavigateListEnum.Quit]: {
        detailTitle: t('app.project.setting.quit.btitle'),
        text: t('app.project.setting.quit.btitle'),
        icon: 'ac-exit',
        sort: 800,
        component: DeleteOrQuit,
        auth: [Authority.Read, Authority.Write],
      },
    })
  }
  navs = filterAuth(navs, auth)

  if (overwrite) {
    Object.keys(navs).forEach((key: any) => {
      const item = navs[key]
      const extendItem = overwrite[key]
      navs[key] = { ...item, ...extendItem }
    })
  }

  return navs
}

function filterModel(navs: any, addon: object) {
  let temp = {} as any
  Object.keys(navs).forEach((key: any) => {
    const val = navs[key]
    if (val.component)
      temp[key] = val
  })
  temp = {
    ...temp,
    ...addon,
  }
  return temp
}

function filterAuth(navs: any, auth: Authority) {
  const temp = {} as any
  Object.keys(navs).forEach((key: any) => {
    const val = navs[key]
    if ((val.auth as string).includes(auth))
      temp[key] = val
  })
  return temp
}

export function getParentNodeKeys(node: Node | undefined) {
  let n = node
  const keys = []
  while (n && n.parent) {
    n = n.parent
    n && n.data?.id && keys.push(n.data?.id)
  }
  return keys
}
