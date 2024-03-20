import type { Nullable } from '../typescript'

const FOCUSABLE_ELEMENT_SELECTORS
  = 'a[href],button:not([disabled]),button:not([hidden]),:not([tabindex="-1"]),input:not([disabled]),input:not([type="hidden"]),select:not([disabled]),textarea:not([disabled])'

/**
 * Determine if the testing element is visible on screen no matter if its on the viewport or not
 */
export function isVisible(element: HTMLElement) {
  const computed = getComputedStyle(element)
  // element.offsetParent won't work on fix positioned
  // potential issue here, going to need some expert advices on this issue
  return computed.position === 'fixed' ? false : element.offsetParent !== null
}

export function obtainAllFocusableElements(element: HTMLElement): HTMLElement[] {
  return Array.from(element.querySelectorAll<HTMLElement>(FOCUSABLE_ELEMENT_SELECTORS)).filter(
    (item: HTMLElement) => isFocusable(item) && isVisible(item),
  )
}

/**
 * @desc Determine if target element is focusable
 * @param element {HTMLElement}
 * @returns true if it is focusable
 */
export function isFocusable(element: HTMLElement): boolean {
  if (element.tabIndex > 0 || (element.tabIndex === 0 && element.getAttribute('tabIndex') !== null))
    return true

  // HTMLButtonElement has disabled
  if ((element as HTMLButtonElement).disabled)
    return false

  switch (element.nodeName) {
    case 'A': {
      // casting current element to Specific HTMLElement in order to be more type precise
      return !!(element as HTMLAnchorElement).href && (element as HTMLAnchorElement).rel !== 'ignore'
    }
    case 'INPUT': {
      return !((element as HTMLInputElement).type === 'hidden' || (element as HTMLInputElement).type === 'file')
    }
    case 'BUTTON':
    case 'SELECT':
    case 'TEXTAREA': {
      return true
    }
    default: {
      return false
    }
  }
}

/**
 * @desc Set Attempt to set focus on the current node.
 * @param element
 *          The node to attempt to focus on.
 * @returns
 *  true if element is focused.
 */
export function attemptFocus(element: HTMLElement): boolean {
  if (!isFocusable(element))
    return false

  // Remove the old try catch block since there will be no error to be thrown
  element.focus?.()
  return document.activeElement === element
}

/**
 * Trigger an event
 * mouseenter, mouseleave, mouseover, keyup, change, click, etc.
 * @param elm
 * @param name
 * @param opts
 */
export const triggerEvent = function (elm: HTMLElement, name: string, ...opts: Array<boolean>): HTMLElement {
  let eventName: string

  if (name.includes('mouse') || name.includes('click'))
    eventName = 'MouseEvents'
  else if (name.includes('key'))
    eventName = 'KeyboardEvent'
  else eventName = 'HTMLEvents'

  const evt = document.createEvent(eventName)

  evt.initEvent(name, ...opts)
  elm.dispatchEvent(evt)
  return elm
}

export const isLeaf = (el: HTMLElement) => !el.getAttribute('aria-owns')

export function getSibling(el: HTMLElement, distance: number, elClass: string): Nullable<Element> {
  const { parentNode } = el
  if (!parentNode)
    return null
  const siblings = parentNode.querySelectorAll(elClass)
  const index = Array.prototype.indexOf.call(siblings, el)
  return siblings[index + distance] || null
}

export function focusNode(el: any) {
  if (!el)
    return
  el.focus()
  !isLeaf(el) && el.click()
}
