export interface SearchFilter {
  jsonPath: number[],
  search: string
}

export const initialSearchFilterState: SearchFilter = {
  jsonPath: [],
  search: ""
}
