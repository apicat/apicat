<template>
  <div ref="domRef" :class="['ac-code-editor', { readonly: readonly }]"></div>
</template>

<script setup lang="ts">
import { indentWithTab } from '@codemirror/commands'
import { EditorView, keymap, highlightSpecialChars, drawSelection, lineNumbers, dropCursor, rectangularSelection, crosshairCursor } from '@codemirror/view'
import { Extension, EditorState } from '@codemirror/state'
import { defaultHighlightStyle, syntaxHighlighting, indentOnInput, bracketMatching, foldGutter } from '@codemirror/language'
import { history, historyKeymap, defaultKeymap } from '@codemirror/commands'
import { autocompletion, closeBrackets } from '@codemirror/autocomplete'
import { useLanguageExtension } from '@/hooks/useCodeMirrorLang'

import { onMounted, onUnmounted, shallowRef, watch } from 'vue'

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

let view: EditorView

const emits = defineEmits(['update:modelValue'])
const { lang } = toRefs(props)
const viewRef: Ref<EditorView | null> = shallowRef(null)
const domRef = shallowRef()

const languageExtension = useLanguageExtension(lang, viewRef)

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
  languageExtension,
  lineNumbers(),
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
  keymap.of([...historyKeymap, ...defaultKeymap, indentWithTab]),
])()

onMounted(() => {
  const exts = [
    basicSetup,
    fixedHeightEditor,
    EditorState.readOnly.of(props.readonly),
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

  viewRef.value = view
})

onUnmounted(() => {
  view?.destroy()
})

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

  .cm-link {
    color: #032f62;
    text-decoration: underline;
  }

  .cm-heading {
    font-weight: bold;
  }

  .cm-emphasis {
    font-style: italic;
  }

  .cm-strong {
    font-weight: bold;
  }

  .cm-keyword {
    color: #708;
  }

  .cm-atom,
  .cm-bool,
  .cm-url,
  .cm-contentSeparator,
  .cm-labelName {
    color: #219;
  }

  .cm-literal,
  .cm-inserted {
    color: #164;
  }

  .cm-deleted {
    color: #b31d28;
    background-color: #ffeef0;
  }

  .cm-inserted {
    color: #22863a;
    background-color: #f0fff4;
  }

  .cm-comment {
    color: #6a737d;
  }

  .cm-operator {
    color: #d73a49;
  }

  .cm-propertyName {
    color: #005cc5;
  }

  .cm-keyword {
    color: #d73a49;
  }

  .cm-string {
    color: #28a745;
  }

  .cm-string2 {
    color: #28a745;
  }

  .cm-typeName {
    color: #005cc5;
  }

  .cm-function.cm-definition {
    color: #6f42c1;
  }

  .cm-variableName.cm-definition {
    color: #e36209;
  }

  .cm-invalid {
    color: #b31d28;
  }

  .cm-number {
    color: #005cc5;
  }
}
</style>
