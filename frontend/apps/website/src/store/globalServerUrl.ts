import { defineStore } from 'pinia'
import {
  apiCreateProjectURL,
  apiDeleteProjectURL,
  apiEditProjectURL,
  apiGetProjectURLList,
  apiSortProjectURLList,
} from '@/api/project/setting/url'

export const useGlobalServerUrlStore = defineStore('project.globalServerUrl', {
  state: () => ({
    urls: [] as ProjectAPI.ResponseURL[],
  }),

  getters: {},

  actions: {
    async getGlobalServerUrlList(projectID: string) {
      this.urls = await apiGetProjectURLList(projectID)
      return this.urls
    },

    async createGlobalServerUrl(projectID: string, data: ProjectAPI.RequestCreateURL) {
      const res = await apiCreateProjectURL(projectID, data)
      this.urls.push(res)
      return res
    },

    async editGlobalServerUrl(projectID: string, urlID: number, data: ProjectAPI.RequestEditURL) {
      const res = await apiEditProjectURL(projectID, urlID, data)
      const url = this.urls.find(val => val.id === urlID)
      if (url) {
        url.url = data.url
        url.description = data.description
      }
      return res
    },

    async deleteGlobalServerUrl(projectID: string, urlID: number) {
      await apiDeleteProjectURL(projectID, urlID)
      this.urls = this.urls.filter((val) => {
        return val.id !== urlID
      })
    },

    async sortGlobalServerUrl(projectID: string, data: ProjectAPI.RequestSortURL) {
      await apiSortProjectURLList(projectID, data)
      const sortlist = data.serverIDs || []
      // 重新排序 .filter(val => val) as ProjectAPI.ResponseURL[]
      this.urls = sortlist.map(id => this.urls.find(url => url.id === id)!)
    },
  },
})
