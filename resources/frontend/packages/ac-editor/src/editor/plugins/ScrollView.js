import { Plugin } from "prosemirror-state";
import Extension from "../lib/Extension";
import scrollIntoView from "smooth-scroll-into-view-if-needed";

const scroll = (view) => {
  if (!view.state.selection.empty) return false
  const dom = view.domAtPos(view.state.selection.$head.start())
  const isScrollTo = dom.node.classList.contains('placeholder')
  if (isScrollTo && dom.node !== view.dom) {
    scrollToElem(dom.node)
  }
}

const scrollToElem = (node) => {
  node && scrollIntoView(node,{
    block: "center",
  })
}

export default class ScrollView extends Extension {
  get name() {
    return "scroll_view";
  }

  get plugins() {
    if (!this.options.enabled){
      return []
    }
    return [
      new Plugin({
        props: {
          handleDOMEvents: {
            keyup: (view) => {
              scroll(view)
            }
          }
        },
      }),
    ];
  }
}
