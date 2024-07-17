import {selectSavefileState} from "./savefile.reducer";
import {createSelector} from "@ngrx/store";
import {flattenContainedObjects, GameCardFormData, SaveFile, trimName, tryCardInit} from "./savefile.state";
import {ObjectState} from "../../types/ttstypes";
import {selectSearchFilter} from "../savefile-search-filters/search-filter.reducer";


export const selectSaveLoadingState = createSelector(
  selectSavefileState, (state: SaveFile) => state.loadingState
);

export const selectSaveName = createSelector(selectSavefileState, trimName);

export const selectSavefileHeader = createSelector(
  selectSavefileState,
  (state: SaveFile) => {
    const data = state.saveData
    if (!data) return;
    return {
      date: data.Date,
      name: data.SaveName,
      description: data.Note,
      path: state.path
    }
  }
);

// TODO: optimise the chain of the selectors, because you are not utilizing memoization correctly.
// TODO2: separate selector post-filter
export const selectFlattenedObjects = createSelector(
  selectSavefileState,
  selectSearchFilter,
  (state, {jsonPath, search}) => {
    const data = state.saveData
    if (!data) return;
    return flattenContainedObjects(data.ObjectStates).filter(object => {
      const pathCondition = jsonPath.length === 0 || object.jsonRelPath.join('/').startsWith(jsonPath.join('/'));
      const searchCondition = search.length === 0 || Object.values(object.object).some(value => {
        return value.toString().replaceAll("[\\][\s\S]\g", "").includes(search);
      });
      return pathCondition && searchCondition;
    });
  }
);

export const selectCollections = createSelector(
  selectFlattenedObjects,
  selectSearchFilter,
  (objects, {jsonPath}) => {
    if (!objects) return;
    return objects
      .filter(({object}) => ['Deck', 'Bag'].includes(object.Name))
      .filter(({jsonRelPath}) => (jsonRelPath.length === jsonPath.length + 1))
  }
);

//TODO cleanup
export const selectCards = createSelector(
  selectFlattenedObjects,
  (objects) => {
    if (!objects) return;
    return objects.filter(({object}) => object.Name.includes('Card'))
  }
);

export const selectCardForms = createSelector(
  selectCards,
  (state?: {
    jsonRelPath: number[];
    object: ObjectState
  }[]): Map<string, GameCardFormData> | undefined => {
    if (!state) {
      return;
    }
    return state.reduce((accum, {jsonRelPath, object}) => {
      const form: GameCardFormData | undefined = tryCardInit(jsonRelPath, object);
      if (form) {
        accum.set(jsonRelPath.join('/'), form);
      }
      return accum;
    }, new Map<string, GameCardFormData>());
  }
)

export const selectObjectPathMap = createSelector(
  selectFlattenedObjects,
  (state) => {
    if (!state) return;
    return state.reduce((map, {object, jsonRelPath}) =>
        map.set(jsonRelPath.join('/'), object),
      new Map<string, ObjectState>());
  }
);

export const selectAllObjectsPathMap = createSelector(
  selectSavefileState, (state) => {
    if (!state.saveData?.ObjectStates) return;
    return flattenContainedObjects(state.saveData.ObjectStates).reduce((map, {object, jsonRelPath}) =>
        map.set(jsonRelPath.join('/'), object),
      new Map<string, ObjectState>());
  })
