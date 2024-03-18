declare module GlobalAPI{
  interface RequestTable {
    page: number
    pageSize: number
  }

  interface ResponseTableData<T> {
    items?: T
    datas?: T
  }

  type ResponseTableDataKey = keyof ResponseTable

  interface ResponseTable<T> extends ResponseTableData<T> {
    count: number
    currentPage: number
    totalPage: number
  }

  type Merge<M, N> = Omit<M, Extract<keyof M, keyof N>> & N
}
