import Extension from '../lib/Extension'
import { NodeSelection, Plugin } from 'prosemirror-state'
import { Decoration, DecorationSet } from 'prosemirror-view'
import { HANDLE_ICON } from '../../common/constants'
// import { getColumnIndex, getRowIndex } from "@/editor/utils";
// import { findParentNodeOfType } from "prosemirror-utils";
import HandleMenu from './HandleMenu'

export default class DragAndDropHandle extends Extension {
    constructor() {
        super()
    }

    get allowCreate() {
        return !this.editor.options.readonly
    }

    bindEditor(editor = null) {
        this.editor = editor
        this.dragHandleMenu = new HandleMenu(editor)
    }

    get name() {
        return 'DragAndDropHandle'
    }

    createDragHandle() {
        const handle = document.createElement('span')
        handle.setAttribute('contenteditable', 'false')
        const icon = document.createElement('span')
        icon.innerHTML = HANDLE_ICON
        handle.appendChild(icon)
        handle.classList.add('handle')
        return handle
    }

    selectionToNode(editorView, event) {
        const position = editorView.posAtCoords({
            left: event.x + 40,
            top: event.y,
        })
        const resolved = editorView.state.doc.resolve(position.pos)

        try {
            const tr = editorView.state.tr
            tr.setSelection(NodeSelection.create(editorView.state.doc, resolved.before(1)))
            editorView.dispatch(tr)
        } catch (e) {
            console.log(e)
        }

        return resolved
    }

    get plugins() {
        return [
            new Plugin({
                props: {
                    decorations: (state) => {
                        const decos = []
                        if (state.doc.childCount >= 1) {
                            state.doc.forEach((node, pos) => {
                                decos.push(
                                    Decoration.widget(pos + 1, this.createDragHandle, {
                                        ignoreSelection: true,
                                        // key:pos
                                    })
                                )

                                decos.push(
                                    Decoration.node(pos, pos + node.nodeSize, {
                                        class: 'draggable',
                                    })
                                )
                            })
                        }
                        return DecorationSet.create(state.doc, decos)
                    },

                    handleDOMEvents: {
                        click: (editorView, event) => {
                            const target = event.target

                            if (target.classList.contains('handle')) {
                                event.stopPropagation()
                                event.preventDefault()

                                const resolved = this.selectionToNode(editorView, event)
                                const selectionDom = editorView.nodeDOM(resolved.before(1))

                                let selectionNode = editorView.state.selection.node

                                this.dragHandleMenu && this.dragHandleMenu.show(target, selectionNode, selectionDom)
                            }
                        },
                        mousedown: (editorView, event) => {
                            const target = event.target
                            if (target.classList.contains('handle')) {
                                this.selectionToNode(editorView, event)
                            }
                        },

                        mouseover: (editorView, event) => {
                            const target = event.target
                            if (!target.classList.contains('handle')) {
                                this.reset()
                            }
                        },
                    },
                },
            }),
        ]
    }

    reset() {
        this.handleMenu && this.handleMenu.destroy()
    }

    destroy() {
        this.reset()

        this.handleMenu = null
    }
}
