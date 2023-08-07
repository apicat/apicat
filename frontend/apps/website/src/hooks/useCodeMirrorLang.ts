import { Compartment, Extension } from '@codemirror/state'
import { languages } from '@codemirror/language-data'
import { EditorView } from '@codemirror/view'
import { LanguageDescription } from '@codemirror/language'

export const useCodeMirrorCompartment = (editorViewRef: Ref<EditorView | null>): [Extension, (extension: Extension) => void] => {
  const compartment = new Compartment()

  const updateCompartment = (extension: Extension) => {
    const dispatch = editorViewRef.value?.dispatch
    if (dispatch) {
      dispatch({
        effects: compartment.reconfigure(extension),
      })
    }
  }

  return [compartment.of([]), updateCompartment]
}

const getCodeMirrorLanguageData = (language: string) => languages.find((lang: LanguageDescription) => lang.name.toLocaleLowerCase() === language.toLocaleLowerCase())

export const useLanguageExtension = (languageRef: Ref<string>, editorViewRef: Ref<EditorView | null>): Extension => {
  const [languageCompartment, updateCompartment] = useCodeMirrorCompartment(editorViewRef)

  async function loadLanguage(language: string) {
    const lang = getCodeMirrorLanguageData(language)
    if (lang) {
      const languageExtension = await lang.load()
      updateCompartment(languageExtension)
    }
  }

  watch(
    languageRef,
    async (newLanguage: string) => {
      if (!newLanguage) {
        return
      }
      await loadLanguage(newLanguage)
    },
    {
      immediate: true,
    }
  )

  return languageCompartment
}
