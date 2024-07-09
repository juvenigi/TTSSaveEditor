import {selectSavefileState} from "./savefile.reducer";
import {createSelector} from "@ngrx/store";
import {SaveFile, trimName} from "./savefile.models";
import {ObjectState} from "../../types/ttstypes";


export const selectSaveLoadingState = createSelector(
  selectSavefileState, (state: SaveFile) => state.loadingState)

export const selectSaveName = createSelector(selectSavefileState, trimName);

export const selectSavefileHeader = createSelector(
  selectSavefileState,
  (state: SaveFile) => {
    const data = state.saveData
    if (!data) return;
    return {
      date: data.Date,
      name: data.SaveName,
      description: data.Note
    }
  }
)

export const selectFlattenedObjects = createSelector(
  selectSavefileState,
  (state: SaveFile) => {
    const data = state.saveData
    if (!data) return;
    return flattenContainedObjects(data.ObjectStates);
  }
);

function flattenContainedObjects(objects: ObjectState[]): {
  objects: ObjectState[],
  parentOf: Map<ObjectState, ObjectState>,
  pathMap: Map<ObjectState, number[]>
} {
  const parentOf: Map<ObjectState, ObjectState> = new Map();
  const pathMap: Map<ObjectState, number[]> = new Map();
  return {
    parentOf,
    pathMap,
    objects: deepVisitContainedObjectsRecur(objects, parentOf, pathMap, []),
  }
}

function deepVisitContainedObjectsRecur(
  objects: ObjectState[],
  parentOf: Map<ObjectState, ObjectState>,
  pathMap: Map<ObjectState, number[]>,
  currentPath: number[]
): ObjectState[] {
  return objects.flatMap((obj: ObjectState, idx: number): ObjectState[] => {
    const finalCurrentPath = [...currentPath, idx];
    pathMap.set(obj, finalCurrentPath)
    if (!obj.ContainedObjects) return [obj];
    obj.ContainedObjects!.forEach((contained: ObjectState) => parentOf.set(contained, obj));
    return [obj, ...deepVisitContainedObjectsRecur(obj.ContainedObjects, parentOf, pathMap, finalCurrentPath)];
  });
}

