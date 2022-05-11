import { render, h } from 'vue'
import tippy from 'tippy.js'
import Feedback from './Feedback.vue'

const vNode = h(Feedback)
const feedbackDom = document.createElement('div')
render(vNode, feedbackDom)

let popper = null

export default {
    show(el) {
        if (!el) {
            return
        }

        if (popper) {
            popper.show()
            return
        }

        popper = tippy(el, {
            placement: 'bottom',
            // trigger: 'click',
            content: vNode.el,
            arrow: true,
            interactive: true,
            theme: 'light',
        })

        popper.show()
    },
}
