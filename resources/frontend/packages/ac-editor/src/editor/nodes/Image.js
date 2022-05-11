import { imagePasteAndDropPlugin, imageUploadPlaceholderPlugin } from '../plugins'
import ImageView from '../extension-nodes/vue-node-view/ImageView.vue'
import Node from '../lib/Node'
import { noop } from 'lodash-es'
import { nodeInputRule } from '../commands'
import { ImageDisplay } from '../../common/constants'

/**
 * Matches following attributes in Markdown-typed image: [, alt, src, class]
 *
 * Example:
 * ![Lorem](image.jpg) -> [, "Lorem", "image.jpg"]
 * ![](image.jpg "class") -> [, "", "image.jpg", "small"]
 * ![Lorem](image.jpg "class") -> [, "Lorem", "image.jpg", "small"]
 */

const IMAGE_INPUT_REGEX = /!\[(.+|:?)]\((\S+)(?:(?:\s+)["'](\S+)["'])?\)/

export default class Image extends Node {
    get defaultOptions() {
        return {
            uploadImage: noop,
            onImageUploadStart: noop,
            onImageUploadStop: noop,
        }
    }

    get name() {
        return 'image'
    }

    get schema() {
        return {
            attrs: {
                src: {
                    default: '',
                },
                alt: {
                    default: '',
                },
                title: {
                    default: '',
                },
                width: {
                    default: null,
                },
                height: {
                    default: null,
                },
                alignment: {
                    default: ImageDisplay.CENTER,
                },
            },
            content: 'inline*',
            atom: true,
            group: 'block',
            isolating: true,
            draggable: true,
            parseDOM: [
                {
                    tag: 'img[src]',
                    getAttrs: (dom) => {
                        let { width, height } = dom.style

                        width = width || dom.getAttribute('width') || null
                        height = height || dom.getAttribute('height') || null

                        const alignment = dom.getAttribute('data-display') || ImageDisplay.CENTER

                        return {
                            src: dom.getAttribute('src'),
                            title: dom.getAttribute('title'),
                            alt: dom.getAttribute('alt'),
                            width: width == null ? null : parseInt(width, 10),
                            height: height == null ? null : parseInt(height, 10),
                            alignment,
                        }
                    },
                },
            ],
            // toDOM: () => ["div"],
            toDOM: (node) => {
                let attrs = { ...node.attrs }
                delete attrs.alignment
                return ['div', { style: node.attrs.alignment ? `text-align: ${node.attrs.alignment}` : null }, ['img', attrs]]
            },
        }
    }

    get component() {
        return ImageView
    }

    toMarkdown(state, node) {
        let markdown = ' ![' + state.esc((node.attrs.alt || '').replace('\n', '') || '') + '](' + state.esc(node.attrs.src)
        if (node.attrs.layoutClass) {
            markdown += ' "' + state.esc(node.attrs.layoutClass) + '"'
        } else if (node.attrs.title) {
            markdown += ' "' + state.esc(node.attrs.title) + '"'
        }
        markdown += ')'
        state.write(markdown)
    }

    parseMarkdown() {
        return {
            node: 'image',
            getAttrs: (token) => {
                return {
                    src: token.attrGet('src'),
                    alt: (token.children[0] && token.children[0].content) || null,
                }
            },
        }
    }

    commands({ type }) {
        return {
            deleteImage: () => (state, dispatch) => {
                dispatch(state.tr.deleteSelection())
                return true
            },

            alignRight: () => (state, dispatch) => {
                const attrs = {
                    ...state.selection.node.attrs,
                    title: null,
                    alignment: ImageDisplay.RIGHT,
                }
                const { selection } = state
                dispatch(state.tr.setNodeMarkup(selection.from, undefined, attrs))
                return true
            },

            alignLeft: () => (state, dispatch) => {
                const attrs = {
                    ...state.selection.node.attrs,
                    title: null,
                    alignment: ImageDisplay.LEFT,
                }
                const { selection } = state
                dispatch(state.tr.setNodeMarkup(selection.from, undefined, attrs))
                return true
            },

            alignCenter: () => (state, dispatch) => {
                const attrs = { ...state.selection.node.attrs, alignment: ImageDisplay.CENTER }
                const { selection } = state
                dispatch(state.tr.setNodeMarkup(selection.from, undefined, attrs))
                return true
            },

            createImage: (attrs) => (state, dispatch) => {
                const { selection } = state
                const position = selection.$cursor ? selection.$cursor.pos : selection.$to.pos
                const node = type.create(attrs)
                const transaction = state.tr.insert(position, node)
                dispatch(transaction)
                return true
            },
        }
    }

    inputRules({ type }) {
        return [
            nodeInputRule(IMAGE_INPUT_REGEX, type, (match) => {
                const [, alt, src, title] = match
                return {
                    src,
                    alt,
                    title,
                }
            }),
        ]
    }

    get plugins() {
        return [imageUploadPlaceholderPlugin, imagePasteAndDropPlugin(this.options)]
    }
}
