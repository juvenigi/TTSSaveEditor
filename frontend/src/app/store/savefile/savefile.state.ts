import {ObjectState, SaveState} from "../../types/ttstypes";
import {generateJSONPatch, Operation} from "generate-json-patch";

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

/**
 * this is the script that every custom card contains in order to render text
 */
export const defaultCardLuaScript: string = "function onload(saved_data)\r\n    --Loads the tracking for if the game has started yet\n    valueData = \"\"\r\n    if saved_data then\r\r\n        valueData = JSON.decode(saved_data)[1]\r\r\n    end\r\n\r\n\r\r\n    self.createInput {\r\n        input_function = \"onInput\",\r\n        function_owner = self,\r\n        label          = \"Write here\",\r\n        alignment      = 3,\r\n        position       = {0,0.4,0},\r\n        rotation       = {0,90,0},\r\r\n        width          = 1300,\r\n        height         = 860,\r\r\n        value          = valueData,\n        font_size      = 80,\r\n    }\r\nend\r\n\r\nfunction onInput(self, ply, text, selected)\r\n    if not selected then\r\n        valueData = text\r\n        updateSave()\r\n    end\r\nend\r\n\r\nfunction updateSave()\r\n    saved_data = JSON.encode({valueData})\r\n    self.script_state = saved_data\r\nend\r\n";

export function setFontSizeOfCustomCard(fontSize: number) {
  return defaultCardLuaScript.replace(/font_size\s*=\s*(\d+)/, `font_size = ${fontSize}`)
}

export function getFontSizeOfCustomCard(obj: ObjectState): number {
  const fontSize = /font_size\s*=\s*(\d+)/;
  return +obj.LuaScript!.match(fontSize)![1]
}

export function objectIsEditableCard(obj: ObjectState): boolean {
  const onLoad = /onLoad/i
  const createInput = /self\.createInput/
  const inputFunc = /"onInput"/i

  return !!obj.LuaScript
    && onLoad.test(obj.LuaScript)
    && createInput.test(obj.LuaScript)
    && inputFunc.test(obj.LuaScript)
}

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

export type GameCardFormData = Exclude<ReturnType<typeof tryCardInit>, undefined>;

export function tryCardInit(jsonRelPath: number[], object: ObjectState) {
  const isCardLike = ['Card', 'CardCustom'].includes(object.Name) && objectIsEditableCard(object);
  if (!isCardLike) return undefined;
  let cardText = '';
  let cardScript;
  let fontSize = -1;
  const path = jsonRelPath.join("/");
  try {
    cardText = object.LuaScriptState ? JSON.parse(object.LuaScriptState ?? '[]')[0] : "";
    cardScript = object.LuaScript;
    fontSize = getFontSizeOfCustomCard(object);
  } catch (_) {
    console.error("cannot parse");
    console.error(object, object.LuaScript, object.LuaScriptState);
    return undefined;
  }
  if (!cardText || !cardScript) return undefined;

  return {cardScript, cardText, path, fontSize};
}

export function getObjectPatch(formData: GameCardFormData, target: ObjectState) {
  try {
    const objPath = formData.path.replaceAll('/','/ContainedObjects/');
    const edit = {
      ...target,
      LuaScriptState: JSON.stringify([formData.cardText]),
      LuaScript: setFontSizeOfCustomCard(formData.fontSize)
    } satisfies ObjectState
    return generateJSONPatch(target, edit)
      .map(operation => ({...operation, path: `/ObjectStates/${objPath}${operation.path}`} satisfies Operation));
  } catch (_) {
    return undefined;
  }
}
