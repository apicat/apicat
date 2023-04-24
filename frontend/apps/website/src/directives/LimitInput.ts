const inputLimit = {
  mounted: (el: any, binding: any) => {
    const arg = binding.arg || 'en'
    // 限制规则
    const restrictRule = {
      zh: /[^\s\u4E00-\u9FA5\s]/g, // 中文
      en: /[^a-zA-Z\s]/g, // 英文
      number: /[^0-9]/g, // 纯数字
    }
    let inputLock = false // 输入锁 使用输入法时关闭限制 结束输入法输入时触发限制
    const doRule = (e: any) => {
      e.target.value = e.target.value.replaceAll((restrictRule as any)[arg], '')
      // 手动更新绑定值
      e.target.dispatchEvent(new Event('input'))
    }
    const target = el instanceof HTMLInputElement ? el : el.querySelector('input')
    target.addEventListener('input', (event: any) => {
      if (!inputLock && event.inputType === 'insertText') {
        doRule(event)
        event.returnValue = false
      }
      event.returnValue = false
    })
    // /* 使用输入法开始触发 */
    target.addEventListener('compositionstart', (event: any) => {
      inputLock = true
    })
    // /* 结束输入法使用触发 */
    target.addEventListener('compositionend', (event: any) => {
      inputLock = false
      doRule(event)
    })
  },
}

const install = function (app: any) {
  app.directive('inputLimit', inputLimit)
}

export default install
