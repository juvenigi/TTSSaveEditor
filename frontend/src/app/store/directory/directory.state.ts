
export interface Directory {
  rootPath: string | undefined
  relPath: string
  directoryEntries: string[]
  loadingState: "DONE" | "LOADING" | "PENDING" | "FAILED"
}

export const initialDirectoryState: Directory = {
  rootPath: undefined,
  relPath: "",
  directoryEntries: [],
  loadingState: "PENDING"
}

export function deduceDividerStyle(state:Directory ) {
  return (state.rootPath ?? '\\').includes('\\')? '\\' : '/';
}
