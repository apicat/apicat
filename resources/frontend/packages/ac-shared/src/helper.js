import { DEFAULT_VAL } from './constant'

export const isNumber = (str) => /^(0|[1-9][0-9]*)$/.test(str)

export function isUrl(url) {
    // return true
    return /^https?:\/\/[\w-.]+(:\d+)?/i.test(url)
}

export function isChinese(v) {
    return /^[^\u4e00-\u9fa5]*$/.test(v)
}

export function isEmail(v) {
    return /^[A-Z0-9._+-]+@[A-Z0-9.-]+\.[A-Z]{2,63}$/i.test(v)
}

// 判断参数是否是其中之一
export function oneOf(value, validList) {
    // eslint-disable-next-line no-plusplus
    for (let i = 0; i < validList.length; i++) {
        if (value === validList[i]) {
            return true
        }
    }
    return false
}

// 驼峰转换
export function camelcaseToHyphen(str) {
    return str.replace(/([a-z])([A-Z])/g, '$1-$2').toLowerCase()
}

// 获取元素相对位置
export function inViewport(element) {
    let actualTop = element.offsetTop
    let current = element.offsetParent
    const elementScrollTop = Math.max(document.body.scrollTop, document.documentElement.scrollTop)

    while (current !== null) {
        actualTop += current.offsetTop
        current = current.offsetParent
    }

    return !(actualTop - elementScrollTop > document.documentElement.clientHeight)
}

// 绑定事件兼容
export const on = (function () {
    if (document.addEventListener) {
        return function (element, event, handler) {
            if (element && event && handler) {
                element.addEventListener(event, handler, false)
            }
        }
    }
    return function (element, event, handler) {
        if (element && event && handler) {
            element.attachEvent(`on${event}`, handler)
        }
    }
})()

// 移除事件兼容
export const off = (function () {
    if (document.removeEventListener) {
        return function (element, event, handler) {
            if (element && event) {
                element.removeEventListener(event, handler, false)
            }
        }
    }
    return function (element, event, handler) {
        if (element && event) {
            element.detachEvent(`on${event}`, handler)
        }
    }
})()

// 获取一个月多少天
export const getMonthOfDay = (month) => {
    const d = new Date()
    month && d.setMonth(month)
    d.setDate(0)
    return d.getDate()
}

// 添加滚动事件
export const addWheelListener = (elem, callback, useCapture) => {
    let prefix = ''
    let _addEventListener
    let support

    // 兼容事件
    if (window.addEventListener) {
        _addEventListener = 'addEventListener'
    } else {
        _addEventListener = 'attachEvent'
        prefix = 'on'
    }

    support =
        'onwheel' in document.createElement('div')
            ? 'wheel' // 各个厂商的高版本浏览器都支持"wheel"
            : document.onmousewheel !== undefined
            ? 'mousewheel' // Webkit 和 IE一定支持"mousewheel"
            : 'DOMMouseScroll' // 低版本firefox

    _addWheelListener(elem, support, callback, useCapture)

    if (support === 'DOMMouseScroll') {
        _addWheelListener(elem, 'MozMousePixelScroll', callback, useCapture)
    }

    function _addWheelListener(elem, eventName, callback, useCapture) {
        elem[_addEventListener](
            prefix + eventName,
            support === 'wheel'
                ? callback
                : (originalEvent) => {
                      !originalEvent && (originalEvent = window.event)

                      const event = {
                          originalEvent,
                          target: originalEvent.target || originalEvent.srcElement,
                          type: 'wheel',
                          deltaMode: originalEvent.type === 'MozMousePixelScroll' ? 0 : 1,
                          deltaX: 0,
                          deltaZ: 0,
                          preventDefault() {
                              originalEvent.preventDefault ? originalEvent.preventDefault() : (originalEvent.returnValue = false)
                          },
                      }

                      if (support === 'mousewheel') {
                          event.deltaY = (-1 / 40) * originalEvent.wheelDelta
                          originalEvent.wheelDeltaX && (event.deltaX = (-1 / 40) * originalEvent.wheelDeltaX)
                      } else {
                          event.deltaY = originalEvent.detail
                      }
                      // 改变 IE 下 this指向问题
                      return callback.call(elem, event)
                  },
            useCapture || false
        )
    }
}

// 判断是否为Trident内核浏览器(IE等)函数
export const browserIsIe = () => !!window.ActiveXObject || 'ActiveXObject' in window

// 下载图片
export const downloadImage = (imgSrc, fileName) => {
    // 创建iframe并赋值的函数,传入参数为图片的src属性值.
    function createIframe(img) {
        const iframe = document.getElementById('IframeReportImg')
        // 如果隐藏的iframe不存在则创建
        if (!iframe) {
            const iframe = document.createElement('iframe')
            iframe.setAttribute('id', 'IframeReportImg')
            iframe.setAttribute('style', 'display: none;')
            iframe.setAttribute('name', 'IframeReportImg')
            iframe.setAttribute('onload', 'downloadImg()')
            iframe.setAttribute('width', '0')
            iframe.setAttribute('height', '0')
            iframe.setAttribute('src', 'about:blank')
            document.body.appendChild(iframe)
        }
        // iframe的src属性如不指向图片地址,则手动修改,加载图片
        if (iframe.getAttribute('src') !== img) {
            iframe.setAttribute('src', img)
        } else {
            // 如指向图片地址,直接调用下载方法
            downloadImg()
        }
    }

    // 下载图片的函数
    function downloadImg() {
        if (document.getElementById('IframeReportImg').getAttribute('src') !== 'about:blank') {
            window.frames.IframeReportImg.document.execCommand('SaveAs')
        }
    }

    function createA(img) {
        const a = document.createElement('a')
        a.setAttribute('id', 'downloadImg')
        a.setAttribute('style', 'display: none;')
        a.setAttribute('href', img)
        a.setAttribute('download', fileName || img)

        !document.getElementById('downloadImg') && document.body.appendChild(a)
        document.getElementById('downloadImg').click()
    }

    browserIsIe() ? createIframe(imgSrc) : createA(imgSrc)
}

// 下载图片
export const downloadFile = (link) => {
    // 创建iframe并赋值的函数,传入参数为图片的src属性值.
    function createIframe(img) {
        const iframe = document.getElementById('IframeReportImg')
        // 如果隐藏的iframe不存在则创建
        if (!iframe) {
            const iframe = document.createElement('iframe')
            iframe.setAttribute('id', 'IframeReportImg')
            iframe.setAttribute('style', 'display: none;')
            iframe.setAttribute('name', 'IframeReportImg')
            iframe.setAttribute('onload', 'downloadImg()')
            iframe.setAttribute('width', '0')
            iframe.setAttribute('height', '0')
            iframe.setAttribute('src', 'about:blank')
            document.body.appendChild(iframe)
        }
        // iframe的src属性如不指向图片地址,则手动修改,加载图片
        if (iframe.getAttribute('src') !== img) {
            iframe.setAttribute('src', img)
        } else {
            // 如指向图片地址,直接调用下载方法
            downloadImg()
        }
    }

    // 下载图片的函数
    function downloadImg() {
        if (document.getElementById('IframeReportImg').getAttribute('src') !== 'about:blank') {
            window.frames.IframeReportImg.document.execCommand('SaveAs')
        }
    }

    function createA(img) {
        const eleLink = document.createElement('a')
        eleLink.setAttribute('id', 'downloadImg')
        eleLink.setAttribute('style', 'display: none;')
        eleLink.setAttribute('href', img)
        document.body.appendChild(eleLink)
        eleLink.click()
        document.body.removeChild(eleLink)
    }

    browserIsIe() ? createIframe(link) : createA(link)
}

// 文字省略
export const ellipsis = (str, len = 10) => {
    if (!str) {
        return DEFAULT_VAL
    }

    const reg = /[\u4e00-\u9fa5]/g
    const slice = str.substring(0, len)
    const chineseCharNum = ~~(slice.match(reg) && slice.match(reg).length)
    const realLen = slice.length * 2 - chineseCharNum
    return str.substr(0, realLen) + (realLen < str.length ? '...' : '')
}

// 设置请求头
export const setHeader = (headers) => {
    headers = headers || {}

    headers['Accept-Version'] = 'v4.0'

    return headers
}

export const bindKey = (arr) => (arr || []).map((item) => ({ ...item, sku: shortid() }))

export const shortid = (prefix = '') => prefix + Math.random().toString(32).substr(3)

export function getPropByPath(obj, path, strict) {
    let tempObj = obj
    path = path.replace(/\[(\w+)\]/g, '.$1')
    path = path.replace(/^\./, '')

    let keyArr = path.split('.')
    let i = 0
    for (let len = keyArr.length; i < len - 1; ++i) {
        if (!tempObj && !strict) break
        let key = keyArr[i]
        if (key in tempObj) {
            tempObj = tempObj[key]
        } else {
            if (strict) {
                throw new Error('please transfer a valid prop path to form item!')
            }
            break
        }
    }
    return {
        o: tempObj,
        k: keyArr[i],
        v: tempObj ? tempObj[keyArr[i]] : null,
    }
}

export function isNaN(val) {
    return Number.isNaN(val)
}

export function clamp(val, min, max) {
    if (val < min) {
        return min
    }
    if (val > max) {
        return max
    }
    return val
}

export function loadCss(src) {
    return new Promise((resolve, reject) => {
        let head = document.getElementsByTagName('head')[0]
        let style = document.createElement('link')

        style.href = src
        style.type = 'text/css'
        style.rel = 'stylesheet'

        style.onload = function () {
            resolve()
        }

        style.onerror = function () {
            reject()
        }

        head.append(style)
    })
}
