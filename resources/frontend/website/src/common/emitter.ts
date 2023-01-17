import mitt from 'mitt'

// mitt 事件集合

export const IS_SHOW_DOCUMENT_TITLE = 'is.show.document.title'

export const DOCUMENT_SAVE_ING = 'doc.save.ing'
export const DOCUMENT_SAVE_DONE = 'doc.save.done'
export const DOCUMENT_SAVE_ERROR = 'doc.save.error'

export const DOCUMENT_SAVE_BTN_ING = 'doc.save.btn.ing'
export const DOCUMENT_SAVE_BTN_DONE = 'doc.save.btn.done'
export const DOCUMENT_SAVE_BTN_ERROR = 'doc.save.btn.error'

export default mitt()
