declare namespace ShareAPI {
  interface ResponseSendShareKey {
    expiration: number
    shareCode: string
  }

  // project
  interface ResponseProjectShareStatus {
    projectID?: string
    hasShare?: boolean
    permission: ProjectAPI.Authority
    visibility: ProjectAPI.Visibility
  }
  type ProjectAuthInfo = ResponseProjectShareStatus
  interface ResponseProjectShareInfo {
    secretKey: string
    permission: ProjectAPI.Authority
    visibility: ProjectAPI.Visibility
  }
  interface ResponseChangeProjectShareStatus {
    secretKey: string
  }
  type ResponseResetProjectShareKey = ResponseChangeProjectShareStatus
  type ResponseCheckProjectShareKey = ResponseSendShareKey

  // document
  interface ResponseDocShareStatus {
    projectID?: string
    collectionID?: number
  }
  export interface ResponseDocShareInfo {
    collectionPublicID: string
    secretKey: string
    visibility: ProjectAPI.Visibility
  }
  export interface ResponseChangeDocShareStatus {
    collectionPublicID: string
    secretKey: string
  }
  export interface ResponseResetDocShareStatus {
    secretKey: string
  }
}
