<template>
  <div :class="['ac-code-editor', { readonly: readonly }]">
    <div class="scroll-content">
      <div class="sticky z-10 text-right top-4px">
        <button class="copy-btn" @click="handlerCopy">
          <el-icon class="mr-2px"><ac-icon-ep-copy-document /></el-icon>
          <span>{{ $t('app.common.copy') }}</span>
        </button>
      </div>
      <div ref="domRef"></div>
    </div>
  </div>
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
import { useCopy } from '@/hooks/useCopy'

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

const handlerCopy = () => useCopy(viewRef.value?.state.doc.toString() || '')

// Event handler for the paste event
const handlePasteEvent = (event: ClipboardEvent, view: EditorView) => {
  try {
    if (lang.value === 'json') {
      const pastedText = event.clipboardData?.getData('text') || ''
      const formatJson = JSON.stringify(JSON.parse(pastedText), null, 2)

      // Replace the original pasted text with the formatted JSON
      const { from, to } = view.state.selection.main
      view.dispatch({
        changes: { from, to, insert: formatJson },
        scrollIntoView: true,
      })

      // Prevent the default paste behavior to avoid duplicating the content
      event.preventDefault()
    }
  } catch (error) {
    //
  }
}

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
    EditorView.domEventHandlers({
      paste: handlePasteEvent,
    }),
    EditorView.lineWrapping
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
  position: relative;

  &.readonly {
    max-height: fit-content;
  }

  &:hover .copy-btn {
    display: flex;
  }

  .copy-btn {
    @apply absolute top-4px right-4px hidden items-center rounded bg-zinc-200 px-6px py-2px text-14px;
  }
}
</style>
