<template>
  <div ref="domRef" :class="['ac-code-editor', { readonly: readonly }]"></div>
</template>

<script setup lang="ts">
import { EditorView } from 'codemirror'
import { Compartment } from '@codemirror/state'
import { indentWithTab } from '@codemirror/commands'
import { json } from '@codemirror/lang-json'
import { xml } from '@codemirror/lang-xml'
import { html } from '@codemirror/lang-html'

import { keymap, highlightSpecialChars, drawSelection, lineNumbers, dropCursor, rectangularSelection, crosshairCursor } from '@codemirror/view'
import { Extension, EditorState } from '@codemirror/state'
import { defaultHighlightStyle, syntaxHighlighting, indentOnInput, bracketMatching, foldGutter, foldKeymap } from '@codemirror/language'
import { defaultKeymap, history, historyKeymap } from '@codemirror/commands'
import { searchKeymap, highlightSelectionMatches } from '@codemirror/search'
import { autocompletion, completionKeymap, closeBrackets, closeBracketsKeymap } from '@codemirror/autocomplete'
import { lintKeymap } from '@codemirror/lint'

import { onMounted, onUnmounted, shallowRef, watch } from 'vue'

interface langType {
  [key: string]: Extension
}
const langs: langType = {
  json: json(),
  xml: xml(),
  html: html(),
  raw: [],
}
let view: EditorView
let language = new Compartment()

const props = withDefaults(
  defineProps<{
    modelValue: string
    lang?: string
    readonly?: boolean
  }>(),
  {
    lang: 'raw',
    modelValue: '',
    readonly: false,
  }
)

const emits = defineEmits(['update:modelValue'])

const domRef = shallowRef()

const fixedHeightEditor = EditorView.baseTheme({
  '&': { fontSize: '14px' },
  '.cm-gutters': {
    backgroundColor: 'transparent',
    borderRight: '1px var(--el-border-color-lighter) solid',
    color: 'var(--el-text-color-disabled)',
  },
  '.cm-scroller': { overflow: 'auto' },
  '.cm-content, .cm-gutter': {
    minHeight: '100px',
    border: 0,
    lineHeight: 1.6,
  },
  '&.cm-editor.cm-focused': {
    outline: 0,
  },
})

const basicSetup: Extension = (() => [
  lineNumbers(),
  // highlightActiveLineGutter(),
  highlightSpecialChars(),
  history(),
  foldGutter(),
  drawSelection(),
  dropCursor(),
  EditorState.allowMultipleSelections.of(true),
  indentOnInput(),
  syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
  bracketMatching(),
  closeBrackets(),
  autocompletion(),
  rectangularSelection(),
  crosshairCursor(),
  highlightSelectionMatches(),
  keymap.of([...closeBracketsKeymap, ...defaultKeymap, ...searchKeymap, ...historyKeymap, ...foldKeymap, ...completionKeymap, ...lintKeymap]),
])()

onMounted(() => {
  const exts = [
    basicSetup,
    fixedHeightEditor,
    keymap.of([indentWithTab]),
    language.of(langs[props.lang]),
    EditorView.editable.of(!props.readonly),
    EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        emits('update:modelValue', view.state.doc.toString())
      }
    }),
  ]
  view = new EditorView({
    doc: props.modelValue || '',
    extensions: exts,
    parent: domRef.value,
  })
})

onUnmounted(() => {
  view?.destroy()
})

onUnmounted(() => {
  view?.destroy()
})

watch(
  () => props.lang,
  () => {
    let langext = langs[props.lang]
    if (!langext) {
      langext = []
    }
    view?.dispatch({
      effects: language.reconfigure(langext),
    })
  }
)

watch(
  () => props.modelValue,
  () => {
    if (props.modelValue == view?.state.doc.toString()) {
      return
    }
    view?.dispatch({
      changes: {
        from: 0,
        to: view.state.doc.length,
        insert: props.modelValue,
      },
    })
  }
)
</script>

<style lang="scss" scoped>
.ac-code-editor {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: var(--el-border-radius-base);
  max-height: 400px;
  overflow-y: scroll;
  background: #f5f5f5;

  &.readonly {
    max-height: fit-content;
  }
}

.ac-code-editor:focus-within {
  box-shadow: rgba(0, 0, 0, 0.1) 0 2px 4px;
}
</style>
