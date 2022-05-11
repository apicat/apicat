import { nextTick } from 'vue'
import javascript from 'highlight.js/lib/languages/javascript'
import sql from 'highlight.js/lib/languages/sql'
import json from 'highlight.js/lib/languages/json'
import xml from 'highlight.js/lib/languages/xml'

export function useHighlight() {
    function initHighlight(els: HTMLElement | NodeList) {
        nextTick(async () => {
            if (!els) return

            const hljs = (await import('highlight.js/lib/core')).default

            hljs.registerLanguage('javascript', javascript)
            hljs.registerLanguage('sql', sql)
            hljs.registerLanguage('json', json)
            hljs.registerLanguage('xml', xml)

            if (els instanceof NodeList)
                els.forEach((el: any) => {
                    hljs.highlightElement(el)
                })
            else hljs.highlightElement(els)
        })
    }

    return { initHighlight }
}
