import type { Ref } from 'vue'
import { unref, computed, onBeforeUnmount } from 'vue'

export const useBackground = (target?: HTMLElement | Ref<HTMLElement> | null, color = '#fff') => {
    const elRef = computed(() => unref(target) || window?.document?.body)
    if (elRef.value?.style) elRef.value.style.setProperty('background', color)

    onBeforeUnmount(() => {
        if (elRef.value?.style) elRef.value.style.removeProperty('background')
    })
}
