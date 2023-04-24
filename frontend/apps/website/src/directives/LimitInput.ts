const inputLimit = {
  mounted: (el: any) => {
    let inputLock = false

    const doRule = (e: any) => {
      e.target.value = e.target.value.match(/[\w-_]+/g, '').join('')
      // 手动更新绑定值
      e.target.dispatchEvent(new Event('input'))
    }

    const target = el instanceof HTMLInputElement ? el : el.querySelector('input')

    el._handler = function (event: any) {
      if (!inputLock && event.inputType === 'insertText') {
        doRule(event)
        event.returnValue = false
      }
      event.returnValue = false
    }

    el._compositionstart = () => {
      inputLock = true
    }

    el._compositionend = (event: any) => {
      inputLock = false
      doRule(event)
    }

    target.addEventListener('input', el._handler)
    target.addEventListener('compositionstart', el._compositionstart)
    target.addEventListener('compositionend', el._compositionend)
  },

  unmounted: (el: any) => {
    const target = el instanceof HTMLInputElement ? el : el.querySelector('input')
    el._handler && target.removeEventListener('input', el._handler)
    el._compositionstart && target.removeEventListener('compositionstart', el._compositionstart)
    el._compositionend && target.removeEventListener('compositionend', el._compositionend)
  },
}

const install = function (app: any) {
  app.directive('inputLimit', inputLimit)
}

export default install
