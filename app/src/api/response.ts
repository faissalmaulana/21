interface RawResponse<T> {
  status: number
  data: T
  error: null | {
    message: string
  }
}

interface Pagination {
  page: number
  size: number
  total_items_in_page: number
  total_items: number
  total_pages: number
}

export type { RawResponse, Pagination }
