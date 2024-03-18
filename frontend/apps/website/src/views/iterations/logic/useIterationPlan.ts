import { traverseTree } from '@apicat/shared'
import type { IterationFormProps } from '../components/IterationForm.vue'
import type { CreateOrEditIteration } from './useIterationForm'
import useApi from '@/hooks/useApi'
import { apiGetCollections } from '@/api/project/collection'
import { CollectionTypeEnum } from '@/commons'

export const defaultProps = {
  label: 'title',
  children: 'items',
  rootValue: 0,
  parentKey: 'parentID',
}

export function useIterationPlan(props: IterationFormProps, iterationInfo: Ref<CreateOrEditIteration>) {
  const collections = ref<CollectionAPI.ResponseCollection[]>([])
  const [isLoadingForTree, getCollections] = useApi(apiGetCollections)
  const selectedCollectionKeys = ref<CollectionAPI.ResponseCollection[]>([])
  const iterationIDRef = toRef(props, 'iterationID')

  // 规划迭代树改变数据后更新form字段
  const onTreeChange = (
    checkedCollections: CollectionAPI.ResponseCollection[],
  ) => {
    iterationInfo.value.collectionIDs = checkedCollections.map(
      item => item.id,
    )
  }

  // 更新默认选中Key值
  function updateSelectedCollectionsForTree() {
    selectedCollectionKeys.value = []
    traverseTree<CollectionAPI.ResponseCollection>(
      (node) => {
        if (node.type !== CollectionTypeEnum.Dir)
          node.selected && selectedCollectionKeys.value.push(node)
        return true
      },
      collections.value,
      { subKey: 'items' },
    )

    iterationInfo.value.collectionIDs = selectedCollectionKeys.value.map(item => item.id)
  }

  // 重新选择项目后重新获取项目下对应的接口文档目录树
  watch(
    () => iterationInfo.value.projectID,
    async (projectID) => {
      if (!projectID) {
        collections.value = []
        return
      }

      collections.value = (await getCollections(
        projectID,
        iterationIDRef.value!,
      )) || []

      updateSelectedCollectionsForTree()
    },
  )

  return {
    defaultProps,
    collections,
    selectedCollectionKeys,
    isLoadingForTree,
    onTreeChange,
  }
}
