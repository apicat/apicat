import Ajax from './Ajax'

import ApicatLogo from '@/assets/image/logo-apicat@2x.png'
// import PdfLogo from '@/assets/image/logo-pdf@2x.png'
import PostmanLogo from '@/assets/image/logo-postman@2x.png'

export const exportDocument = (doc = {}) => Ajax.post('/api_doc/export', doc)
export const getExportDocumentResult = (project_id, job_id) => Ajax.get('/api_doc/export_result', { params: { project_id, job_id } })

export const API_PROJECT_EXPORT_ACTION_MAPPING = [
    { text: 'ApiCat', icon: ApicatLogo, type: 'apicat', action: exportDocument, getJobResult: getExportDocumentResult },
    // { text: 'PDF', icon: PdfLogo, type: 'pdf', action: exportDocument, getJobResult: getExportDocumentResult },
    { text: 'Postman(v2.1)', icon: PostmanLogo, type: 'postman', action: exportDocument, getJobResult: getExportDocumentResult },
]

export const API_SINGLE_EXPORT_ACTION_MAPPING = API_PROJECT_EXPORT_ACTION_MAPPING
