import {selectSavefileState} from "./savefile.reducer";
import {createSelector} from "@ngrx/store";
import {SaveFile, trimName} from "./savefile.models";
import {SaveState} from "../../types/ttstypes";


export const selectSaveLoadingState = createSelector(
  selectSavefileState, (state: SaveFile) => state.loadingState)
export const selectSaveName = createSelector(selectSavefileState, trimName);

export const selectSaveRaw = selectSavefileState;
export const selectTabletopSave = createSelector(
  selectSavefileState, (state: SaveFile) => {
    //FIXME: validate the structure
    return state.saveData as unknown as SaveState | null;
  });

