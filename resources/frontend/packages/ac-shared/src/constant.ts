// 默认值
export const DEFAULT_VAL = '--'

// 默认图片集合
export const DEFAULT_IMAGE = {}

// http methods
export const HTTP_METHODS = {
    TYPES: [
        { text: 'GET', value: 1, color: '#66BE74' },
        { text: 'POST', value: 2, color: '#4894FF' },
        { text: 'PUT', value: 3, color: '#51B9C3' },
        { text: 'PATCH', value: 4, color: '#F1924E' },
        { text: 'DELETE', value: 5, color: '#DF4545' },
        { text: 'OPTION', value: 6, color: '#A973DF' },
    ],

    valueOf(value: number) {
        return this.TYPES.find((item) => item.value === value)
    },
}

// 请求体数据格式
export const REQUEST_BODY_DATA_TYPES = {
    TYPES: [
        { text: 'none', value: 0, tip: '无参数' },
        { text: 'form-data', value: 1, tip: '文本、文件' },
        { text: 'x-www-form-urlencoded', value: 2, tip: '纯文本' },
        { text: 'raw', value: 3, tip: '任意格式的文本，text、json、xml、html等等' },
        { text: 'binary', value: 4, tip: '二进制数据(文件)' },
    ],
    valueOf(value: number) {
        return (this.TYPES.find((item) => item.value === value) || { text: DEFAULT_VAL }).text
    },
}

// mock数据类型
export const PARAM_TYPES = {
    TYPES: [
        { text: 'Int', value: 1 },
        { text: 'Float', value: 2 },
        { text: 'String', value: 3 },
        { text: 'Array', value: 4 },
        { text: 'Object', value: 5 },
        { text: 'ArrayObject', value: 8 },
        { text: 'Boolean', value: 6 },
        { text: 'File', value: 7 },
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

// mock.js 数据类型
export const MOCK_DATA_TYPES = [
    { name: '字符串', mock: '@string' },
    { name: '自然数', mock: '@natural' },
    { name: '浮点数', mock: '@float' },
    { name: '字符', mock: '@character' },
    { name: '布尔', mock: '@boolean' },
    { name: 'url', mock: '@url' },
    { name: '域名', mock: '@domain' },
    { name: 'ip地址', mock: '@ip' },
    { name: 'id', mock: '@id' },
    { name: 'guid', mock: '@guid' },
    { name: '当前时间', mock: '@now' },
    { name: '时间戳', mock: '@timestamp' },
    { name: '日期', mock: '@date' },
    { name: '时间', mock: '@time' },
    { name: '日期时间', mock: '@datetime' },
    { name: '图片连接', mock: '@image' },
    { name: '图片data', mock: '@imageData' },
    { name: '颜色', mock: '@color' },
    { name: '颜色hex', mock: '@hex' },
    { name: '颜色rgba', mock: '@rgba' },
    { name: '颜色rgb', mock: '@rgb' },
    { name: '颜色hsl', mock: '@hsl' },
    { name: '整数', mock: '@integer' },
    { name: 'email', mock: '@email' },
    { name: '大段文本', mock: '@paragraph' },
    { name: '句子', mock: '@sentence' },
    { name: '单词', mock: '@word' },
    { name: '大段中文文本', mock: '@cparagraph' },
    { name: '中文标题', mock: '@ctitle' },
    { name: '标题', mock: '@title' },
    { name: '姓名', mock: '@name' },
    { name: '中文姓名', mock: '@cname' },
    { name: '中文姓', mock: '@cfirst' },
    { name: '中文名', mock: '@clast' },
    { name: '英文姓', mock: '@first' },
    { name: '英文名', mock: '@last' },
    { name: '中文句子', mock: '@csentence' },
    { name: '中文词组', mock: '@cword' },
    { name: '地址', mock: '@region' },
    { name: '省份', mock: '@province' },
    { name: '城市', mock: '@city' },
    { name: '地区', mock: '@county' },
    { name: '转换为大写', mock: '@upper' },
    { name: '转换为小写', mock: '@lower' },
    { name: '挑选（枚举）', mock: '@pick' },
    { name: '打乱数组', mock: '@shuffle' },
    { name: '协议', mock: '@protocol' },
]

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
