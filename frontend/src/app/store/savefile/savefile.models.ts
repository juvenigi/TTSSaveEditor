import {SaveState} from "../../types/ttstypes";

export interface SaveFile {
  path: string
  saveData?: SaveState
  objectEdits: Map<string, any>
  loadingState: "DONE" | "LOADING" | "PENDING" | "FAILED"
}

export const initialSaveFileState: SaveFile = {
  path: "",
  saveData: undefined,
  objectEdits: new Map(),
  loadingState: "PENDING"
};

export function trimName(save: SaveFile) {
  const lastSlash = save.path.lastIndexOf('/')
  return save.path.substring(lastSlash + 1);
}
