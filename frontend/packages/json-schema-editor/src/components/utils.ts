export function typename(type: string | string[] | undefined) {
  if (type === undefined) {
    return 'any'
  }
  if (type instanceof Array) {
    return type.length > 1 ? 'other' : type[0]
  }
  return type
}

export function setValueByPath(obj: any, path: string, value: any) {
  const arr = path.split('.')
  let current = obj
  for (let i = 0; i < arr.length - 1; i++) {
    current = current[arr[i]]
  }
  current[arr[arr.length - 1]] = value
}
