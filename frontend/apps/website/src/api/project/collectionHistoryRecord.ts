import dayjs from 'dayjs'
import { parseJSONWithDefault } from '@apicat/shared'
import Ajax from '../Ajax'

const baseRestfulApiPath = (projectID: string, collectionID: number) => `/projects/${projectID}/collections/${collectionID}/histories`

// 获取集合历史记录列表
export function getCollectionHistoryRecords(projectID: string, collectionID: number): Promise<HistoryRecord.ResponseCollectionRecordList> {
  return Ajax.get(baseRestfulApiPath(projectID, collectionID), { params: { startTime: dayjs().subtract(3, 'month').unix(), endTime: dayjs().unix() } })
}

// 获取集合历史记录详情
export async function getCollectionHistoryRecordDetail(projectID: string, collectionID: number, historyID: number): Promise<HistoryRecord.CollectionDetail> {
  const info: HistoryRecord.CollectionDetail = await Ajax.get(`${baseRestfulApiPath(projectID, collectionID)}/${historyID}`)
  info.content = parseJSONWithDefault(info.content, [])
  return info
}

// 恢复集合历史
export function restoreCollectionHistoryRecord(projectID: string, collectionID: number, historyID: number): Promise<void> {
  return Ajax.put(`${baseRestfulApiPath(projectID, collectionID)}/${historyID}/restore`)
}

// 集合历史记录对比
export function compareCollection(projectID: string, collectionID: number, originalID: number, targetID: number): Promise<HistoryRecord.CollectionDiff> {
  return Ajax.get(`${baseRestfulApiPath(projectID, collectionID)}/diff`, { params: { originalID, targetID } })
}
