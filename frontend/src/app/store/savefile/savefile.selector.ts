import {selectSavefileState} from "./savefile.reducer";
import {createSelector} from "@ngrx/store";
import {flattenContainedObjects, SaveFile, trimName, tryCardInit} from "./savefile.state";
import {selectSearchFilter} from "../search-filters/search-filter.reducer";
import {FormControl} from "@angular/forms";
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
  selectSearchFilter,
  (state, {jsonPath, search}) => {
    const data = state.saveData
    if (!data) return;
    return flattenContainedObjects(data.ObjectStates).filter(object => {
      const pathCondition = jsonPath.length === 0 || object.jsonRelPath.join('/') === jsonPath.join('/');
      const searchCondition = search.length === 0 || Object.values(object.object).some(value => {
        return value.toString().replaceAll("[\\][\s\S]\g", "").includes(search);
      });
      return pathCondition && searchCondition;
    });
  }
);

export const selectCardForms = createSelector(
  selectFlattenedObjects,
  (state?: { jsonRelPath: number[]; object: ObjectState }[]): Map<string, Exclude<ReturnType<typeof tryCardInit>, undefined>> | undefined => {
    if (!state) return;
    return state.reduce((accum, obj) => {
      const form: undefined | FormControl<string> = tryCardInit(obj.object);
      if (form) {
        accum.set(obj.jsonRelPath.join('/'), form);
      }
      return accum;
    }, new Map<string, FormControl>);
  }
)




