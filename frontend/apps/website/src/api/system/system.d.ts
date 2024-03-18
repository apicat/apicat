declare namespace SystemAPI {
  type StorageDrivers = import('@/commons/constant').SysStorage
  type CacheDrivers = import('@/commons/constant').SysCache
  type EmailDrivers = import('@/commons/constant').SysEmail
  type ModelDrivers = import('@/commons/constant').SysModel

  // service
  interface ServiceData {
    appName: string
    appServerBind: string
    appUrl: string
    mockServerBind: string
    mockUrl: string
  }

  // oauth
  interface OAuthData {
    clientID: string
    clientSecret: string
  }

  // storage
  interface StorageDisk {
    path: string
  }
  interface StorageCF {
    accountID: string
    accessKeyID: string
    accessKeySecret: string
    bucketName: string
    bucketUrl: string
  }
  interface StorageQiniu {
    accessKey: string
    secretKey: string
    bucketName: string
    bucketUrl: string
  }
  interface StorageItem {
    driver: StorageDrivers
    use: boolean
    config: StorageDisk | StorageCF | StorageQiniu
  }

  // cache
  interface CacheMemory {}
  interface CacheRedis {
    host: string
    password?: string
    database: number
  }
  interface CacheItem {
    driver: CacheDrivers
    use: boolean
    config: CacheMemory | CacheRedis
  }

  // database
  interface Database {
    host: string
    username: string
    // password: string
    database: string
  }

  // email
  interface EmailSMTP {
    host: string
    user?: string
    address: string
    password: string
  }
  interface EmailSendCloud {
    apiUser: string
    apiKey: string
    fromEmail: string
    fromName: string
  }
  interface EmailItem {
    driver: EmailDrivers
    use: boolean
    config: EmailSMTP | EmailSendCloud
  }

  // model
  interface ModelAzure {
    apiKey: string
    endpoint: string
    llmName: string
  }
  interface ModelOpenAI {
    apiKey: string
    organizationID?: string
    apiBase?: string
    llmName: string
  }
  interface ModelItem {
    driver: ModelDrivers
    use: boolean
    config: ModelAzure | ModelOpenAI
  }
}
