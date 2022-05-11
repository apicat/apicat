import { EditorView } from 'prosemirror-view'
import { Schema } from 'prosemirror-model'
import { dropCursor } from 'prosemirror-dropcursor'
import { inputRules } from 'prosemirror-inputrules'
import { keymap } from 'prosemirror-keymap'
import { AllSelection, EditorState, NodeSelection, TextSelection } from 'prosemirror-state'
import { baseKeymap, chainCommands, liftEmptyBlock, newlineInCode, splitBlock } from 'prosemirror-commands'

import ComponentView from './lib/ComponentView'
import Emitter from './lib/Emitter'
import { gapCursor } from './plugins/prosemirror-gapcursor'
import ExtensionManager from './lib/ExtensionManager'
import dictionary from './dictionary'

import {
    Blockquote,
    BulletList,
    CheckboxItem,
    CheckboxList,
    Doc,
    Embed,
    HardBreak,
    Heading,
    HorizontalRule,
    Image,
    ListItem,
    OrderedList,
    Paragraph,
    Table,
    TableCell,
    TableHeadCell,
    TableRow,
    Text,
} from './nodes'

import ApiUrl from './extension-nodes/ApiUrl'
import HttpUrl from './extension-nodes/HttpUrl'
import HttpCode from './extension-nodes/HttpCode'
import ParamTable from './extension-nodes/ParamTable'
import HttpRequestParam from './extension-nodes/HttpRequestParam'
import HttpResponseParam from './extension-nodes/HttpResponseParam'

import CodeBlockHighlight from './extension-nodes/CodeBlock'
import Bold from './marks/Bold'
import Code from './marks/Code'
import Italic from './marks/Italic'
import Underline from './marks/Underline'
import Link from './marks/Link'
import Highlight from './marks/Highlight'
import Strike from './marks/Strike'

import { noop } from 'lodash-es'

import { BlockMenuTrigger, DragAndDropHandle, History, Keys, MarkdownPaste, TrailingNode } from './plugins'

import BlockMenu from './plugins/EditorTriggerMenu/BlockMenu'
import FloatingToolbar from './plugins/FloatingToolbar'
import LinkToolbar from './plugins/LinkToolbar'
import MockRules from './plugins/MockRules'
import ScrollView from './plugins/ScrollView'

import NodeEditViewManager from './lib/NodeEditViewManager'
import CommonParamsManager from './lib/CommonParamsManager'
import CommonUrlManager from './lib/CommonUrlManager'
import checkContent from './utils/checkContent'
import checkAllNode from './utils/checkAllNode'

export const EDITOR_EVENTS = {
    Init: 'onInit',
    Transaction: 'onTransaction',
    Update: 'onUpdate',
    Focus: 'onFocus',
    Blur: 'onBlur',
    Paste: 'onPaste',
    Drop: 'onDrop',
    ShowToast: 'onShowToast',
    ClickLink: 'onClickLink',
    ClickHashTag: 'onClickHashTag',
    HoverLink: 'onHoverLink',
    BlockMenuOpen: 'onBlockMenuOpen',
    BlockMenuClose: 'onBlockMenuClose',
    EditResponseMock: 'onEditResponseMock',
}

class AcEditor extends Emitter {
    constructor(el = null, options) {
        super()
        this.element = el
        this.init(options)
    }

    static DEFAULT_OPTIONS = {
        dictionary,
        readonly: false,
        isShowDevTool: false,
        useBaseExtensions: true,
        extensions: [],
        content: '',
        topNode: 'doc',
        emptyDocument: {
            type: 'doc',
            content: [
                {
                    type: 'paragraph',
                },
            ],
        },
        parseOptions: {},
        uploadImage: null,
        onImageUploadStart: null,
        onImageUploadStop: null,
        getAllCommonParams: null,
        addCommonParam: null,
        deleteCommonParam: null,
        showToast: () => {},
        openNotification: () => {},
    }

    static create(el, opt) {
        return new AcEditor(el, opt)
    }

    get state() {
        return this.view ? this.view.state : null
    }

    get isDestroyed() {
        return this.view ? !this.view.docView : !this.view
    }

    get isEditable() {
        return this.view && this.view.editable
    }

    init(options) {
        this.mergeOptions(options)

        this.commonParamsManager = new CommonParamsManager(this)
        this.commonUrlManager = new CommonUrlManager(this)
        this.mockModel = new MockRules(this)

        this.extensions = this.createExtensions()
        this.nodes = this.createNodes()
        this.marks = this.createMarks()
        this.schema = this.createSchema()
        this.plugins = this.createPlugins()
        this.keymaps = this.createKeymaps()
        this.serializer = this.createSerializer()
        this.parser = this.createParser()
        this.inputRules = this.createInputRules()
        this.nodeViews = this.createNodeViews()
        this.view = this.createView()
        this.commands = this.createCommands()

        // extend editor plugin
        this.nodeEditViewManager = this.createNodeEditViewManager()
        this.blockMenu = new BlockMenu(this, this.options)
        this.floatingToolbar = new FloatingToolbar(this, this.options)
        this.linkToolbar = new LinkToolbar(this, this.options)
    }

    mergeOptions(options) {
        this.options = {
            ...AcEditor.DEFAULT_OPTIONS,
            ...options,
        }
    }

    baseExtensions() {
        return [
            new Doc(),
            new Text(),
            new Paragraph(),

            new HardBreak(),
            new Blockquote(),
            new CodeBlockHighlight({ readonly: this.options.readonly }),
            new CheckboxList(),
            new CheckboxItem(),
            new BulletList(),
            new Embed(),
            new ListItem(),
            new Heading({
                dictionary,
            }),
            new HorizontalRule(),
            new Table(),
            new TableCell(),
            new TableHeadCell(),
            new TableRow(),
            new Bold(),
            new Code(),
            new Highlight(),
            new Italic(),
            new Underline(),
            new Strike(),
            new Link({
                onKeyboardShortcut: noop,
                onClickLink: () => this[EDITOR_EVENTS.ClickLink](),
                onClickHashTag: () => this[EDITOR_EVENTS.ClickHashTag](),
                onHoverLink: () => this[EDITOR_EVENTS.HoverLink](),
            }),
            new Strike(),
            new OrderedList(),
            new Image({
                dictionary,
                uploadImage: this.options.uploadImage,
                onImageUploadStart: this.options.onImageUploadStart,
                onImageUploadStop: this.options.onImageUploadStop,
                showToast: this.options.showToast,
            }),

            //api doc node
            new ApiUrl(),
            new HttpUrl(),
            new HttpCode(),
            new ParamTable(),
            new HttpRequestParam(),
            new HttpResponseParam(),

            // plugins
            new History(),
            new Keys({
                onBlur: noop,
                onFocus: noop,
                onSave: noop,
                onSaveAndExit: noop,
                onCancel: noop,
            }),
            new TrailingNode(),
            new MarkdownPaste(),
            new DragAndDropHandle(),
            new BlockMenuTrigger(this.options),
            new ScrollView({ enabled: true }),
        ]
    }

    createExtensions() {
        return new ExtensionManager([...this.baseExtensions(), ...this.options.extensions], this)
    }

    createNodes() {
        return this.extensions.nodes
    }

    createMarks() {
        return this.extensions.marks
    }

    createSchema() {
        return new Schema({
            topNode: this.options.topNode,
            nodes: this.nodes,
            marks: this.marks,
        })
    }

    createPlugins() {
        return this.extensions.plugins
    }

    createKeymaps() {
        return this.extensions.keymaps({
            schema: this.schema,
        })
    }

    createInputRules() {
        return this.extensions.inputRules({
            schema: this.schema,
        })
    }

    createNodeViews() {
        let views = this.extensions.extensions
            .filter((extension) => extension.component)
            .reduce((nodeViews, extension) => {
                const nodeView = (node, view, getPos, decorations, innerDecorations) => {
                    return new ComponentView(extension.component, {
                        editor: this,
                        extension,
                        node,
                        view,
                        getPos,
                        decorations,
                        innerDecorations,
                    })
                }

                return {
                    ...nodeViews,
                    [extension.name]: nodeView,
                }
            }, {})

        views = this.extensions.extensions
            .filter((extension) => extension.nodeView)
            .reduce((nodeViews, extension) => {
                const nodeView = (node, view, getPos, decorations, innerDecorations) => {
                    return new extension.nodeView({
                        extension,
                        node,
                        view,
                        getPos,
                        decorations,
                        innerDecorations,
                        editor: this,
                    })
                }

                return {
                    ...nodeViews,
                    [extension.name]: nodeView,
                }
            }, views || {})

        return views
    }

    createView() {
        if (!this.element) {
            throw new Error('createView called before ref available')
        }
        let view = new EditorView(this.element, {
            state: this.createState(),
            nodeViews: this.nodeViews,
            editable: () => !this.options.readonly,
            dispatchTransaction: (transaction) => {
                let isUpdateContent = transaction.getMeta('update.content')
                if (this.options.readonly && !isUpdateContent && transaction.curSelection instanceof NodeSelection) {
                    return
                }
                const { state } = this.view.state.applyTransaction(transaction)
                this.view.updateState(state)

                if (!transaction.docChanged) {
                    return
                }
                this[EDITOR_EVENTS.Update]()
            },
        })

        if (this.options.readonly) {
            view.dom.className += ' readonly'
        }

        return view
    }

    createState() {
        // copy prosemirror-commands createParagraphNear
        function defaultBlockAt(match) {
            for (var i = 0; i < match.edgeCount; i++) {
                var ref = match.edge(i)
                var type = ref.type
                if (type.isTextblock && !type.hasRequiredAttrs()) {
                    return type
                }
            }
            return null
        }

        function createParagraphNear(state, dispatch) {
            var sel = state.selection
            var $from = sel.$from
            var $to = sel.$to
            if (sel instanceof AllSelection || $from.parent.inlineContent || $to.parent.inlineContent) {
                return false
            }
            var type = defaultBlockAt($to.parent.contentMatchAt($to.indexAfter()))
            if (!type || !type.isTextblock) {
                return false
            }
            if (dispatch) {
                var side = (!$from.parentOffset && $to.index() < $to.parent.childCount ? $from : $to).pos

                if (!$from.parentOffset) {
                    side = $to.pos
                }

                var tr = state.tr.insert(side, type.createAndFill())
                tr.setSelection(TextSelection.create(tr.doc, side + 1))
                dispatch(tr.scrollIntoView())
            }
            return true
        }

        // 覆写默认回车行为
        baseKeymap.Enter = chainCommands(newlineInCode, createParagraphNear, liftEmptyBlock, splitBlock)

        let doc = this.createDocument(this.options.content)
        return EditorState.create({
            schema: this.schema,
            doc,
            plugins: [
                ...this.plugins,
                ...this.keymaps,
                dropCursor({ class: 'drop-cursor' }),
                gapCursor(),
                inputRules({
                    rules: this.inputRules,
                }),
                keymap(baseKeymap),
            ],
        })
    }

    createDocument(content) {
        if (!content) {
            return this.schema.nodeFromJSON(this.options.emptyDocument)
        }

        if (typeof content === 'object') {
            try {
                const node = this.schema.nodeFromJSON(checkContent(content, this.schema))
                return checkAllNode(node)
            } catch (error) {
                return this.schema.nodeFromJSON(this.options.emptyDocument)
            }
        }
    }

    createSerializer() {
        return this.extensions.serializer()
    }

    createParser() {
        return this.extensions.parser({
            schema: this.schema,
        })
    }

    createCommands() {
        return this.extensions.commands({
            schema: this.schema,
            view: this.view,
        })
    }

    createNodeEditViewManager() {
        return new NodeEditViewManager(this.extensions.extensions, this)
    }

    /**
     * Register a ProseMirror plugin.
     */
    registerPlugin(plugin, handlePlugins) {
        const plugins = typeof handlePlugins === 'function' ? handlePlugins(plugin, this.state.plugins) : [...this.state.plugins, plugin]

        const state = this.state.reconfigure({ plugins })

        this.view.updateState(state)
    }

    /**
     * Unregister a ProseMirror plugin.
     */
    unregisterPlugin(nameOrPluginKey) {
        if (this.isDestroyed) {
            return
        }

        const name = typeof nameOrPluginKey === 'string' ? `${nameOrPluginKey}$` : nameOrPluginKey.key

        const state = this.state.reconfigure({
            plugins: this.state.plugins.filter((plugin) => !plugin.key.startsWith(name)),
        })

        this.view.updateState(state)
    }

    getJSON() {
        return this.state.doc.toJSON()
    }

    setContent(content = {}) {
        const { doc, tr } = this.state
        const document = this.createDocument(content)
        const selection = TextSelection.create(doc, 0, doc.content.size)
        const transaction = tr.setSelection(selection).setMeta('update.content', true).replaceSelectionWith(document, false)

        this.view.dispatch(transaction)
    }

    clearContent() {
        this.setContent(this.options.emptyDocument)
    }

    focus() {
        this.view && this.view.focus()
    }

    destroy() {
        this.emit('destroy')
        this.view && this.view.destroy()
        this.nodeEditViewManager && this.nodeEditViewManager.destroy()
        this.blockMenu && this.blockMenu.destroy()
        this.linkToolbar && this.linkToolbar.destroy()
        this.mockModel && this.mockModel.destroy()

        this.blockMenu = null
        this.linkToolbar = null
        this.floatingToolbar = null
        this.mockModel = null
    }
}

Object.keys(EDITOR_EVENTS).forEach(function (eventKey) {
    let method = EDITOR_EVENTS[eventKey]
    AcEditor.prototype[method] = function (args) {
        this.emit(method, this, args)
    }
})

export const createEditor = (el, opt) => AcEditor.create(el, opt)
export default AcEditor
