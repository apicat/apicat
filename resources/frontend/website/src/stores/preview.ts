import { getProjectInfo, getSingleApiDocumentDetail } from '@/api/preview'
import { defineStore } from 'pinia'
import { router } from '@/router'
import { Storage } from '@natosoft/shared'
import { PREVIEW_PROJECT, PREVIEW_PROJECT_SECRET, PREVIEW_DOCUMENT, PREVIEW_DOCUMENT_SECRET } from '@/router/constant'

interface ProjectState {
    projectInfo: { has_shared: boolean; visibility: number; in_this: boolean; [key: string]: any } | null
    documentInfo: { [key: string]: any } | null
}

export const usePreviewStore = defineStore({
    id: 'preview',

    state: (): ProjectState => ({
        projectInfo: null,
        documentInfo: null,
    }),

    getters: {},

    actions: {
        async getProjectInfo(project_id: any) {
            const { data } = await getProjectInfo(project_id)
            this.projectInfo = data || null
            return this.projectInfo
        },

        async getDocumentInfo(doc_id: any) {
            const token = Storage.get(Storage.KEYS.SECRET_DOCUMENT_TOKEN + doc_id, true) || ''
            const { data } = await getSingleApiDocumentDetail(token, doc_id)
            this.documentInfo = data || null
            return this.documentInfo
        },

        // 访问预览页面秘钥失效后
        goVerificationPage() {
            const project_id = router.currentRoute.value.params.project_id || ''
            const document_id = router.currentRoute.value.params.doc_id || ''
            const name = (router.currentRoute.value.name || '') as string

            if (name.includes(PREVIEW_PROJECT)) {
                Storage.remove(Storage.KEYS.SECRET_PROJECT_TOKEN + project_id, true)
                router.push({ name: PREVIEW_PROJECT_SECRET, params: { project_id } })
                return
            }

            if (name.includes(PREVIEW_DOCUMENT)) {
                Storage.remove(Storage.KEYS.SECRET_DOCUMENT_TOKEN + document_id, true)
                router.push({ name: PREVIEW_DOCUMENT_SECRET, params: { document_id } })
            }
        },
    },
})
