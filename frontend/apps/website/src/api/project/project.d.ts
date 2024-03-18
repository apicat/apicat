declare namespace ProjectAPI {
  type Visibility = import('@/commons/constant').Visibility
  type Authority = import('@/commons/constant').Authority

  interface RequestProject {
    isFollowed?: boolean
    groupID?: number
    permissions?: Authority | Authority[] | string
  }

  interface ResponseProject {
    id: string
    title: string
    cover: string | ProjectCover
    memberID: number
    shareKey: string
    createdAt: string
    updatedAt: string
    visibility: Visibility
    selfMember: {
      isFollowed: boolean
      permission: Authority
      groupID?: number
    }
    description?: string
    mockURL?: string
  }

  interface ProjectCover {
    type: 'icon' | 'url'
    coverBgColor?: string
    coverIcon?: string
    coverUrl?: string
  }

  interface RequestCreateProject {
    teamID: string
    cover: string
    title: string
    visibility: Visibility
    data?: string
    description?: string
    groupID?: number
    type?: string
  }

  interface RequestChangeGroup {
    groupID?: number
  }

  interface RequestSetProjectGeneral {
    title: string
    visibility: Visibility
    cover: string
    description?: string
  }

  interface RequestCreateGroup {
    name: string
    teamID: string
  }

  interface ResponseGroup {
    id: string | number
    name: string
  }

  interface RequestRenameGroup {
    groupID?: number
    name: string
  }
  interface RequestSortGroup {
    groupIDs: number[]
  }

  type Status = import('@/commons/constant').Status
  interface RequestCreateMember {
    id: string
    memberIDs: number[]
    permission: Authority
  }

  interface Member {
    id: number
    name: string
    email: string
    permission: Authority
    createdAt: string
    role: Role
    updatedAt: string
    status?: Status
  }
  interface RequestChangeMember {
    id: string
    memberID: number
    permission: Authority
  }
  interface RequestEditMember {
    permission: Authority
  }

  interface GlobalParameterTypes {
    Cookie: 'cookie'
    Header: 'header'
    Path: 'path'
    Query: 'query'
  }

  type GlobalParameterType = GlobalParameterTypes[keyof GlobalParameterTypes]

  interface GlobalParameter {
    in: GlobalParameterType
    name: string
    id: string
    required: boolean
    schema: Record<string, any>
  }

  type ResponseGlobalParamList = {
    [Key in GlobalParameterType]: GlobalParameter[]
  }

  interface RequestCreateURL {
    url: string
    description?: string
  }
  interface ResponseURL {
    id: number
    url: string
    description?: string
  }
  interface RequestEditURL {
    url: string
    description?: string
  }
  interface RequestSortURL {
    serverIDs: number[]
  }

  // Trash
  interface Trash {
    collectionID: number
    collectionTitle: string
    deletedAt: string
    deletedBy: string
  }

  interface RequestRestoreTrash {
    collectionIDs?: number[]
  }

  // test case
  interface TestCase {
    id: number
    title: string
    createdAt: number
  }
  interface TestCaseDetail {
    id: string
    title: string
    content: string
  }
  interface ResponseTestCase {
    generating: boolean
    records: TestCase[]
  }

  // ai
  interface RequestCreateCollectionWithAI {
    parentID: number
    prompt: string
    iterationID?: string
  }
}
