import type { FormInstance } from 'element-plus'
import { debounce } from 'lodash-es'
import useApi from '@/hooks/useApi'
import useProjectStore from '@/store/project'
import { useGlobalServerUrlStore } from '@/store/globalServerUrl'

interface URLItem {
  $rowID: number
  data: ProjectAPI.ResponseURL
}

export function useServerUrlCompare() {
  const projectStore = useProjectStore()
  const serverUrlStore = useGlobalServerUrlStore()
  const [isLoadingUrls, getURLs] = useApi(serverUrlStore.getGlobalServerUrlList)

  // 检测每个URL里的内容是否变更，包括 位置/url/description/新建 (上述都由`isSame(id)`方法检测)
  // 在样式上，若有变更则会在url的input框最前方添加一个圆形的icon标识

  // oldURL用于储存原数据，与newURL进行对比查看是否存在变更
  // urlList为合并oldURL与newURL，newURL覆盖oldURL用于将所有更改都集中在newURL上
  // draglist 单独用于排序, odraglist 对比查看是否存在顺序改变 {uuid,id,url:ref(),description:ref()}
  const oldURL = ref<Record<number, URLItem>>({})
  const newURL = ref<Record<number, URLItem>>({})
  const urlList = computed(() => {
    return {
      ...oldURL.value,
      ...newURL.value,
    }
  })
  const odraglist = ref<number[]>([])
  const draglist = ref<number[]>([])
  let currentID = 0

  // 该方法用户获取url列表并处理，oldURL直接就是原数据，newURL为复制后数据，draglist为取id后的数组
  async function initURL() {
    // 复制以取消和store中数据使用同一份引用
    const res: ProjectAPI.ResponseURL[] = JSON.parse(JSON.stringify((await getURLs(projectStore.project!.id)) || []))
    const d: Record<number, URLItem> = {}
    const l = []
    for (let i = 0; i < res.length; i++) {
      const data = {
        $rowID: currentID++,
        data: res[i],
      }
      d[data.$rowID] = data
      l.push(data.$rowID)
    }
    odraglist.value = l
    draglist.value = JSON.parse(JSON.stringify(l))
    oldURL.value = d
    newURL.value = JSON.parse(JSON.stringify(d))
  }

  // 根据提供的id判断是否存在变化，主要是位置/url/description判断
  // 新添加的url直接返回false
  function isSame(rowID: number, withPosition = false): boolean {
    const o = oldURL.value[rowID]
    const n = newURL.value[rowID]
    if (!(o && n))
      return false
    if (withPosition) {
      function positionSame() {
        if (!(draglist.value.includes(rowID) && odraglist.value.includes(rowID)))
          return true
        else return draglist.value.indexOf(rowID) === odraglist.value.indexOf(rowID)
      }
      if (!positionSame())
        return false
    }
    for (const key in o.data) {
      const ov = (o as any).data[key]
      const nv = (n as any).data[key]
      if (ov !== nv)
        return false
    }
    return true
  }

  // 由于正常的id都大于0，所以本地新添加但未提交的url的id从负数起开始设置
  // 这里的id本身并无太大作用但主要为唯一标识和新建标识所以尽量不要删和改动
  function addURL() {
    const data: URLItem = {
      $rowID: currentID++,
      data: {
        id: -1,
        url: '',
        description: '',
      },
    }
    newURL.value[data.$rowID] = data
    draglist.value.push(data.$rowID)
  }

  // 防抖上传
  const navFormRef = shallowRef<FormInstance>()

  // 防抖
  const saveURL = debounce(async () => {
    // 校验
    const validID: string[] = []
    try {
      await navFormRef.value!.validate((_, validData) => {
        for (const key in validData)
          validID.push(key.split('.')[0])
      })
    }
    catch (e) {}
    let doSort = false
    const task: Promise<any>[] = [] // 存放所有新建和修改url的Promise
    for (const _rowID in newURL.value) {
      if (validID.includes(_rowID))
        continue
      const rowID: number = Number.parseInt(_rowID)
      const val = newURL.value[rowID]
      // 如果id为负数就新建
      if (val.data.id < 0) {
        task.push(
          serverUrlStore
            .createGlobalServerUrl(projectStore.project!.id, {
              url: val.data.url,
              description: val.data.description,
            })
            .then((res) => {
              if (!doSort)
                doSort = true
              newURL.value[rowID].data = res
              oldURL.value[rowID] = JSON.parse(JSON.stringify(newURL.value[rowID]))
            }),
        )
      }
      else if (!isSame(rowID)) {
        // 存在url或description不同就发送编辑请求
        task.push(
          serverUrlStore
            .editGlobalServerUrl(projectStore.project!.id, val.data.id, {
              url: val.data.url,
              description: val.data.description,
            })
            .then(() => {
              oldURL.value[rowID].data.url = val.data.url
              oldURL.value[rowID].data.description = val.data.description
            }),
        )
      }
    }
    // 等待所有新建和创建请求以及后续处理之后再看排序
    await Promise.all(task)

    // 只要新建或存在位置变化就请求排序
    if (doSort || isPositionChanged()) {
      const rowIDs: number[] = []
      const ids: number[] = []
      draglist.value.forEach((val) => {
        const id = newURL.value[val].data.id
        if (id >= 0) {
          ids.push(id)
          rowIDs.push(val)
        }
      })
      await serverUrlStore
        .sortGlobalServerUrl(projectStore.project!.id, {
          serverIDs: ids,
        })
        .then(() => {
          odraglist.value = JSON.parse(JSON.stringify(rowIDs))
        })
    }
  }, 200)

  function isPositionChanged() {
    // // 数量不一样
    // if (odraglist.value.length !== draglist.value.length) return true

    // 位置不一样
    let ii = 0
    for (let i = 0; i < draglist.value.length; i++) {
      const n = draglist.value[i]
      if (newURL.value[n].data.id < 0)
        continue
      const o = odraglist.value[ii]
      if (o !== n)
        return true
      ii++
    }
    return false
  }

  // 删除
  async function deleteURL(rowID: number) {
    const id = newURL.value[rowID].data.id
    if (id >= 0)
      await serverUrlStore.deleteGlobalServerUrl(projectStore.project!.id, id)

    delete oldURL.value[rowID]
    delete newURL.value[rowID]
    draglist.value = draglist.value.filter((item) => {
      return item !== rowID
    })
    odraglist.value = odraglist.value.filter((item) => {
      return item !== rowID
    })
    saveURL()
  }

  initURL()

  return {
    initURL,
    isSame,
    addURL,
    isPositionChanged,
    deleteURL,
    urlList,
    isLoadingUrls,
    draglist,
    saveURL,

    navFormRef,

    // useless
    newURL,
    oldURL,
    odraglist,
  }
}
