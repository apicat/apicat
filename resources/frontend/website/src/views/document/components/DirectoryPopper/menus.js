import { h, render } from 'vue'
import NewMenus from './NewMenus.vue'
import DirOperateMenus from './DirOperateMenus.vue'
import DocOperateMenus from './DocOperateMenus.vue'

import createDocIcon from '@/assets/image/doc-common@2x.png'
import createHttpDocIcon from '@/assets/image/doc-http@2x.png'

export const NEW_MENUS = [
    { text: '新建文档', value: 'createDoc', img: createDocIcon, onClick: 'createDoc' },
    { text: '新建 HTTP API 文档', value: 'createHttpDoc', img: createHttpDocIcon, onClick: 'createHttpDoc' },
    { text: '导入文件', value: 'importFile', icon: 'iconimport', onClick: 'onImportBtnClick', divided: true },
    { text: '新建分类', value: 'newCatetory', type: 'newCatetory', icon: 'iconIconPopoverTree', onClick: 'onCreateDirBtnClick', divided: true },
]

export const DOC_OPERATE_MENUS = [
    { text: '复制', value: 'copy', onClick: 'onCopyBtnClick' },
    { text: '重命名', value: 'rename', onClick: 'onRenameBtnClick' },
    { text: '分享', value: 'share', onClick: 'onShareBtnClick' },
    { text: '导出', value: 'export', onClick: 'onExportBtnClick' },
    { text: '删除', value: 'delete', onClick: 'onDeleteBtnClick' },
]

export const DIR_OPERATE_MENUS = [DOC_OPERATE_MENUS[1], DOC_OPERATE_MENUS[4]]

export const DIR_NEW_TYPE = 'DIR_NEW_TYPE'
export const DIR_OPERATE_TYPE = 'DIR_OPERATE_TYPE'
export const DOC_OPERATE_TYPE = 'DOC_OPERATE_TYPE'

function renderComp(comp, props) {
    const domWrapper = document.createElement('div')
    const vm = h(comp)
    render(vm, domWrapper)

    Object.keys(props).forEach((key) => {
        vm.component.props[key] = props[key]
    })

    return vm
}

export default {
    [DIR_NEW_TYPE]: {
        render(props) {
            return renderComp(NewMenus, props)
        },
    },

    [DIR_OPERATE_TYPE]: {
        render(props) {
            return renderComp(DirOperateMenus, props)
        },
    },

    [DOC_OPERATE_TYPE]: {
        render(props) {
            return renderComp(DocOperateMenus, props)
        },
    },
}
