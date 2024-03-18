import { isClient } from '@vueuse/core'
import type { CSSProperties } from 'vue'
import { isNumber, isString } from '../types'
import { camelize } from '../strings'

export const classNameToArray = (cls = '') => cls.split(' ').filter(item => !!item.trim())

export function hasClass(el: Element, cls: string): boolean {
  if (!el || !cls)
    return false
  if (cls.includes(' '))
    throw new Error('className should not contain space.')
  return el.classList.contains(cls)
}

export function addClass(el: Element, cls: string) {
  if (!el || !cls.trim())
    return
  el.classList.add(...classNameToArray(cls))
}

export function removeClass(el: Element, cls: string) {
  if (!el || !cls.trim())
    return
  el.classList.remove(...classNameToArray(cls))
}

export function getStyle(element: HTMLElement, styleName: keyof CSSProperties): string {
  if (!isClient || !element || !styleName)
    return ''

  let key = camelize(styleName)
  if (key === 'float')
    key = 'cssFloat'
  try {
    const style = (element.style as any)[key]
    if (style)
      return style
    const computed: any = document.defaultView?.getComputedStyle(element, '')
    return computed ? computed[key] : ''
  }
  catch {
    return (element.style as any)[key]
  }
}

export function addUnit(value?: string | number, defaultUnit = 'px') {
  if (!value)
    return ''
  if (isString(value))
    return value
  else if (isNumber(value))
    return `${value}${defaultUnit}`
}
