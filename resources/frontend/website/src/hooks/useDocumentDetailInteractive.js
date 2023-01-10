import { nextTick } from 'vue'
import { getAttr, hasClass, showOrHide, toggleClass } from '@natosoft/shared'
import mediumZoom from 'medium-zoom'
import tippy from 'tippy.js'
import { useHighlight } from '@/hooks/useHighlight'

const expand = (pid, isExpand) => {
    document.querySelectorAll('[data-pid="' + pid + '"]').forEach(function (el) {
        let arrow = el.querySelector('i.editor-arrow-right')
        if (arrow && !hasClass(arrow, 'expand')) {
            toggleClass(arrow, 'expand')
        }
        let id = getAttr(el, 'data-id')
        el.style.display = isExpand ? null : 'none'
        id && expand(id, isExpand)
    })
}

const initTableToggle = (rootSelector) => {
    document.querySelectorAll(`${rootSelector}.ac-param-table .editor-arrow-right`).forEach(function (el) {
        el.onclick = function () {
            expand(getAttr(this, 'data-id'), !hasClass(this, 'expand'))
            toggleClass(this, 'expand')
        }
    })

    document.querySelectorAll(`${rootSelector}div.collapse-title .response_body_title`).forEach(function (el) {
        el.onclick = function () {
            const h3 = this.parentElement
            const parent = h3.parentElement
            const isShow = hasClass(parent, 'close')
            showOrHide(h3.nextElementSibling, isShow)
            showOrHide(parent.nextElementSibling, isShow)
            toggleClass(parent, 'close')
        }
    })

    document.querySelectorAll(`${rootSelector}h3.collapse-title >span`).forEach(function (el) {
        el.onclick = function () {
            const parent = this.parentElement
            const isShow = hasClass(parent, 'close')
            showOrHide(parent.nextElementSibling, isShow)
            toggleClass(parent, 'close')
        }
    })
}

const initMediumZoom = (rootSelector) => {
    mediumZoom(`${rootSelector}.ProseMirror .image-view img`, {
        template: '#template-zoom-image',
        container: '[data-zoom-container]',
    })
}

const initTippy = (rootSelector) => {
    tippy('[data-tippy-content]', { theme: 'light', appendTo: document.querySelector(`${rootSelector}.ProseMirror`) })
}

const initCodeBlockToClipboard = (rootSelector) => {
    document.querySelectorAll(`${rootSelector}.code-block button`).forEach((el) => {
        el.setAttribute('data-text', el.parentElement.querySelector('code').innerText)
    })
}

export const useDocumentDetailInteractive = async (rootSelector) => {
    const { initHighlight } = useHighlight()
    await nextTick()

    const rootSelectorStr = rootSelector ? `${rootSelector} ` : ''

    initTableToggle(rootSelectorStr)
    initTippy(rootSelectorStr)
    initMediumZoom(rootSelectorStr)
    initCodeBlockToClipboard(rootSelectorStr)
    initHighlight(document.querySelectorAll(`${rootSelectorStr}pre code`))
}
