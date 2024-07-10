import {ObjectState, SaveState} from "../../types/ttstypes";
import {FormControl} from "@angular/forms";

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

export function flattenContainedObjects(objects: ObjectState[]) {
  return locateObjectRecur(objects, [])
}

/**
 *
 * @param objects
 * @param currentPath
 * @return flat array of ObjectState with their json paths (array indices)
 */
function locateObjectRecur(
  objects: ObjectState[],
  currentPath: number[]
) {
  return objects.flatMap((obj: ObjectState, idx: number): { jsonRelPath: number[], object: ObjectState }[] => {
    const jsonRelPath = [...currentPath, idx];
    const wrappedObject = {jsonRelPath, object: obj}
    if (!obj.ContainedObjects) {
      return [wrappedObject];
    } else {
      return [wrappedObject, ...locateObjectRecur(obj.ContainedObjects, jsonRelPath)];
    }
  });
}

export function tryCardInit(object: ObjectState) {
  const isCardLike = ['Card', 'CardCustom'].includes(object.Name);
  if (!isCardLike) return undefined;
  let cardText = '';
  try {
    cardText = JSON.parse(object.LuaScriptState ?? '[]')[0];
  } catch (_) {

  }
  if (!cardText) return undefined;
  return new FormControl<string>(cardText, {nonNullable: true});
}
