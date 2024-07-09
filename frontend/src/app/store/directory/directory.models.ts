
export interface Directory {
  rootPath: string | undefined
  directoryEntries: string[]
  loadingState: "DONE" | "LOADING" | "PENDING" | "FAILED"
}

export const initialDirectoryState: Directory = {
  rootPath: undefined,
  directoryEntries: [],
  loadingState: "PENDING"
}
