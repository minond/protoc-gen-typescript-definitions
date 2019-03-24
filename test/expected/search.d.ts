export type SearchRequest = {
  query?: string
  page_number?: number
  result_per_page?: number
  corpus?: number
}

export type SearchResponse = {
  results?: Array<{
    url?: string
    title?: string
    snippets?: string[]
  }>
}
