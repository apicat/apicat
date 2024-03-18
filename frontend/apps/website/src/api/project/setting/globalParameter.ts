import { parseJSONWithDefault } from '@apicat/shared'
import Ajax from '../../Ajax'
import { gatherSharedTokenWithParams } from '@/api/shareToken'

// 获取全局参数列表
export async function getGlobalParameters(projectID: string): Promise<ProjectAPI.ResponseGlobalParamList> {
  const parameters: ProjectAPI.ResponseGlobalParamList = await Ajax.get(`/projects/${projectID}/global/parameters`, { params: gatherSharedTokenWithParams({}, projectID) }, { isShowErrorMsg: false })
  Object.keys(parameters).forEach((key) => {
    const data = (parameters[key as ProjectAPI.GlobalParameterType] = parameters[key as ProjectAPI.GlobalParameterType] || [])
    data.map((item: ProjectAPI.GlobalParameter) => {
      item.schema = parseJSONWithDefault(item.schema, {})
      return item
    })
  })
  return parameters
}

// 新增全局参数
export async function addGlobalParameter(
  projectID: string,
  data: Omit<ProjectAPI.GlobalParameter, 'id'>,
): Promise<ProjectAPI.GlobalParameter> {
  const parameter: ProjectAPI.GlobalParameter = await Ajax.post(`/projects/${projectID}/global/parameters`, {
    ...data,
    schema: JSON.stringify(data.schema),
  })
  parameter.schema = parseJSONWithDefault(parameter.schema, {})
  return parameter
}

// 更新全局参数
export function updateGlobalParameter(projectID: string, { id, ...data }: ProjectAPI.GlobalParameter): Promise<void> {
  return Ajax.put(`/projects/${projectID}/global/parameters/${id}`, { ...data, schema: JSON.stringify(data.schema) })
}

// 删除全局参数
export function deleteGlobalParameter(projectID: string, parameterID: string, deref: boolean): Promise<void> {
  return Ajax.delete(`/projects/${projectID}/global/parameters/${parameterID}?deref=${deref}`)
}

// 排序全局参数
export function sortGlobalParameter(projectID: string, data: { parameterIDs: string[];in: ProjectAPI.GlobalParameterType }): Promise<void> {
  return Ajax.put(`/projects/${projectID}/global/parameters/sort`, data)
}
