import {createFeatureSelector, createSelector} from "@ngrx/store";
import {AppState, reducerFeatureKey} from "./app.reducer";

// Feature selector for the AppState
export const selectAppState = createFeatureSelector<AppState>(reducerFeatureKey);

// Selector for saveData
export const selectSaveData = createSelector(
  selectAppState, (state: AppState) => state?.saveData
);

// Selector for fsPath
export const selectFsPath = createSelector(
  selectAppState, (state: AppState) => state?.fsPath
);
