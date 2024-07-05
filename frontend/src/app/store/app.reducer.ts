import {createFeatureSelector, createReducer, createSelector, on} from '@ngrx/store';
import {AppActions, SaveFileApiActions} from "./app.actions";

export const reducerFeatureKey = 'reducer';

export interface AppState {
  fsPath: string
  savefilePaths: string[]
  saveData: string | undefined
  loadingState: "DONE" | "LOADING" | "PENDING" | "FAILED"
}

export const initialState: AppState = {saveData: "", savefilePaths: [], fsPath: "initial", loadingState: "PENDING"};

export const reducer = createReducer(
  initialState,
  on(AppActions.requestSavefile, (state, {fsPath}) => {
    return {...state, loadingState: (state.fsPath !== fsPath) ? "LOADING" : "DONE"} satisfies AppState;
  }),
  on(SaveFileApiActions.savefileRetrieveSuccess, (state, {data, fsPath}) => {
    console.debug("reducer fired: api success")
    return {...state, loadingState: "DONE", saveData: data, fsPath} satisfies AppState;
  }),
  on(SaveFileApiActions.savefileRetrieveFailure, (_, {message, statusCode, fsPath}) => {
    console.debug(`fetching failed ${statusCode}, ${fsPath}, ${message}`);
    return {...initialState, loadingState: "FAILED"} satisfies AppState;
  }),
  on(SaveFileApiActions.requestSavefilesSuccess, (state, res) => {

    console.debug(res)
    return {...state, savefilePaths: res.paths, fsPath: res.fsPath} satisfies AppState;
  })
);

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

export const selectAvailablePaths = createSelector(
  selectAppState, (state) => state.savefilePaths
)
