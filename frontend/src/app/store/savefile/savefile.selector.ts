import {selectSavefileState} from "./savefile.reducer";
import {createSelector} from "@ngrx/store";
import {flattenContainedObjects, GameCardFormControl, SaveFile, trimName, tryCardInit} from "./savefile.state";
import {FormControl} from "@angular/forms";
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

export const selectFlattenedObjects = createSelector(
  selectSavefileState,
  selectSearchFilter,
  (state, {jsonPath, search}) => {
    const data = state.saveData
    if (!data) return;
    console.debug(jsonPath)
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
  }[]): Map<string, Exclude<ReturnType<typeof tryCardInit>, undefined>> | undefined => {
    if (!state) return;
    return state.reduce((accum, obj) => {
      const form: undefined | FormControl<string> = tryCardInit(obj.object);
      if (form) {
        accum.set(obj.jsonRelPath.join('/'), form);
      }
      return accum;
    }, new Map<string, GameCardFormControl>);
  }
)

export const selectObjectPathMap = createSelector(
  selectFlattenedObjects,
  (state) => {
    if (!state) return;
    state.reduce((map, {object, jsonRelPath}) =>
        map.set(jsonRelPath.join('/'), object),
      new Map<string, ObjectState>);
  }
);




