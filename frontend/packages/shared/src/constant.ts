// 默认值
export const DEFAULT_VAL = '--'

// 默认图片集合
export const DEFAULT_IMAGE = {}

// http methods
export const HTTP_METHODS = {
  TYPES: [
    { text: 'GET', value: 'get', color: '#66BE74' },
    { text: 'POST', value: 'post', color: '#4894FF' },
    { text: 'PUT', value: 'put', color: '#51B9C3' },
    { text: 'PATCH', value: 'patch', color: '#F1924E' },
    { text: 'DELETE', value: 'delete', color: '#DF4545' },
    { text: 'OPTION', value: 'option', color: '#A973DF' },
  ],

  valueOf(value: number) {
    return this.TYPES.find((item) => item.value === value)
  },
}

// 请求体数据格式
export const REQUEST_BODY_DATA_TYPES = {
  TYPES: [
    { text: 'none', value: 'none', tip: '无参数' },
    { text: 'form-data', value: 'HTTP_METHODS', tip: '文本、文件' },
    { text: 'x-www-form-urlencoded', value: 'x-www-form-urlencoded', tip: '纯文本' },
    { text: 'raw', value: 'raw', tip: '任意格式的文本，text、json、xml、html等等' },
    { text: 'binary', value: 'binary', tip: '二进制数据(文件)' },
  ],
  valueOf(value: number) {
    return (this.TYPES.find((item) => item.value === value) || { text: DEFAULT_VAL }).text
  },
}

export const REQUEST_DATA_TYPES = {
  TYPES: [
    { text: 'text', value: 0 },
    { text: 'file', value: 1 },
  ],

  valueOf(value: number) {
    return (this.TYPES.find((item) => item.value === value) || { text: DEFAULT_VAL }).text
  },
}

// 事件KEY_MAP
export const EVENT_CODE = {
  tab: 'Tab',
  enter: 'Enter',
  space: 'Space',
  left: 'ArrowLeft', // 37
  up: 'ArrowUp', // 38
  right: 'ArrowRight', // 39
  down: 'ArrowDown', // 40
  esc: 'Escape',
  delete: 'Delete',
  backspace: 'Backspace',
  numpadEnter: 'NumpadEnter',
  pageUp: 'PageUp',
  pageDown: 'PageDown',
  home: 'Home',
  end: 'End',
}
