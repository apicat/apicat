import { delay } from '@apicat/shared'
import useApi from '@/hooks/useApi'
import { apiCreateTestCase, apiGetTestCaseList, apiRegenTestCaseList } from '@/api/project/collection'

const wait = async () => await delay(5000)

class TestCaseState {
  // button state
  public btnLoading
  private startBtn() {
    this.btnLoading
  }
  private stopBtn() {
    this.btnLoading
  }

  // list and api
  private getTestCasesWithoutLoading = apiGetTestCaseList
  public testCases: ProjectAPI.TestCase[]
  public listLoading: boolean
  public collectionID: string
  public projectID: string
  constructor(projectID: string, collectionID: string) {
    this.running = true
    this.btnLoading = this.listLoading = false
    this.testCases = []
    const [_, getTestCases] = useApi(apiGetTestCaseList)
    this.collectionID = collectionID
    this.projectID = projectID
    getTestCases(projectID, collectionID)
      .then(async (res) => {
        if (!this.running) return
        console.log('update')
        this.testCases = res!.records
        if (res!.generating) this.startLoop()
        this.startBtn()
        await wait()
        this.stopBtn()
      })
      .finally(() => (this.listLoading = false))
  }

  public stopLoop: undefined | (() => void)
  public running: boolean = false
  startLoop(immediate = false) {
    let loopRun = true
    if (this.stopLoop) this.stopLoop()
    this.stopLoop = function () {
      loopRun = false
      this.stopBtn()
      this.stopLoop = undefined
    }
    const loop = async () => {
      this.startBtn()
      if (!immediate) await wait()
      try {
        while (loopRun && this.running) {
          const res = await this.getTestCasesWithoutLoading(this.projectID, this.collectionID)
          this.testCases = res!.records
          if (!res!.generating) break

          await wait()
        }
      } finally {
        this.stopBtn()
      }
    }
    return loop()
  }

  public generateTestCases(prompt?: string) {
    this.startBtn()
    apiCreateTestCase(this.projectID, this.collectionID, prompt)
      .then(() => {
        if (!this.running) return
        this.startLoop()
      })
      .catch(this.stopBtn)
  }

  public reGenerateTestCases() {
    this.startBtn()
    apiRegenTestCaseList(this.projectID, this.collectionID, undefined)
      .then(() => {
        if (!this.running) return
        this.startLoop(true)
      })
      .catch(this.stopBtn)
  }

  public onDel(id: number) {
    this.testCases = this.testCases.filter((val) => val.id !== id)
  }

  public dispose() {
    this.running = false
    if (this.stopLoop) this.stopLoop()
  }
}

export function useTestCaseList(props: { projectID: string; collectionID: string }) {
  const testCaseState = ref<{ t: TestCaseState | undefined }>({ t: undefined })
  // const testCaseState = reactive<{ value: TestCaseState | undefined }>({ value: undefined })
  // const testCases = computed(() => {
  //   console.log(computed)
  //   return testCaseState.value?.testCases ?? []
  // })
  // const listLoading = computed(() => testCaseState.value?.listLoading ?? false)
  // const btnLoading = computed(() => testCaseState.value?.btnLoading ?? false)
  const inited = ref(false)
  const disposeState = () => {
    if (testCaseState.value.t) {
      testCaseState.value.t.dispose()
      testCaseState.value.t = undefined
    }
  }

  function endLoop() {
    disposeState()
  }

  function init() {
    disposeState()
    testCaseState.value.t = reactive(new TestCaseState(props.projectID, props.collectionID))
    inited.value = false
  }
  // const generateTestCases = (prompt?: string) => testCaseState.value?.generateTestCases(prompt)
  // const reGenerateTestCases = () => testCaseState.value?.reGenerateTestCases()
  // const onDel = (id: number) => testCaseState.value?.onDel(id)

  return {
    testCaseState,
    // testCases,
    // btnLoading,
    // listLoading,
    init,
    inited,
    // generateTestCases,
    // reGenerateTestCases,
    // onDel,

    endLoop,
  }
}
