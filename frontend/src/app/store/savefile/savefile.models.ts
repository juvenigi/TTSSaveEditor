import {JSONValue} from "../../misc";

export interface SaveFile {
  path: string
  saveData: JSONValue
  objectEdits: Map<string, any>
  loadingState: "DONE" | "LOADING" | "PENDING" | "FAILED"
}

export const initialSaveFileState: SaveFile = {
  path: "",
  saveData: null,
  objectEdits: new Map(),
  loadingState: "PENDING"
};

export function trimName(save: SaveFile) {
  const lastSlash = save.path.lastIndexOf('/')
  return save.path.substring(lastSlash + 1);
}
