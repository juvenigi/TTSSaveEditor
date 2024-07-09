export interface SeachFilter {
  jsonPath: number[],
  search: string
}

export const initialSearchFilterState: SeachFilter = {
  jsonPath: [],
  search: ""
}
