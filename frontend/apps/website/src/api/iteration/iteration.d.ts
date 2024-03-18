declare namespace IterationAPI {
  interface RequestCreateIteration {
    projectID: string
    title: string
    collectionIDs?: number[]
    description?: string
  }

  interface ResponseIteration {
    id: string
    apisCount?: number
    createdAt: number
    createdBy?: string
    description?: string
    project?: ProjectAPI.ResponseProject
    title: string
    updatedAt: string
    updatedBy?: string
  }

  interface RequestEditIteration {
    title: string
    collectionIDs?: number[]
    description?: string
  }
}
