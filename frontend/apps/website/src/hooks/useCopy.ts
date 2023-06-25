function feakCopy(value: string) {
  const pre = document.createElement('pre')
  const code = document.createElement('code')
  pre.style.position = 'absolute'
  pre.style.opacity = '0'
  pre.className = 'copy_text'
  code.innerText = value
  pre.appendChild(code)
  document.body.appendChild(pre)
  pre.click()
  pre.remove()
}

export const useCopy = (text: string) => feakCopy(text)
