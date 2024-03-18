declare module CollectionAPI {
  type CollectionTypeEnum = import('@/commons/constant').CollectionTypeEnum

  interface RequestCreateCollection {
    title: string
    content?: string
    iterationID?: string
    parentID: number
    projectID: string
    type: CollectionTypeEnum
  }

  interface ResponseCollectionDetail {
    id: number
    parentID: number
    title: string
    type: CollectionTypeEnum
    content?: Array<any>
    publicID?: string
    sharePassword?: string
  }

  interface ResponseCollection {
    id: number
    title: string
    type: CollectionTypeEnum
    items?: ResponseCollection[]
    parentID: number
    selected?: boolean
  }

  interface RequestEditCollectionDetail {
    title: string
    content?: string
  }

  interface RequestMoveCollection {
    origin: CollectionOrder
    target: CollectionOrder
  }

  interface CollectionOrder {
    ids: number[]
    parentID: number
  }
}
