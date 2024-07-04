import {createReducer, on} from '@ngrx/store';
import {AppActions, SaveFileApiActions} from "./app.actions";

export const reducerFeatureKey = 'reducer';

export interface AppState {
  fsPath: string
  saveData: string | undefined
  loadingState: "DONE" | "LOADING" | "PENDING" | "FAILED"
}

export const initialState: AppState = {saveData: "", fsPath: "initial", loadingState: "PENDING"};

export const reducer = createReducer(
  initialState,
  on(AppActions.requestSavefile, (state, {fsPath}) => {
    console.debug("reducer fired")
    return {...state, loadingState: (state.fsPath !== fsPath) ? "LOADING" : "DONE"} satisfies AppState;
  }),
  on(SaveFileApiActions.savefileRetrieveSuccess, (state, {data, fsPath}) => {
    console.debug("reducer fired: api success")
    return {...state, loadingState: "DONE", saveData: data, fsPath} satisfies AppState;
  }),
  on(SaveFileApiActions.savefileRetrieveFailure, (_, {message, statusCode, fsPath}) => {
    console.debug(`fetching failed ${statusCode}, ${fsPath}, ${message}`);
    return {...initialState, loadingState: "FAILED"} satisfies AppState;
  })
);

