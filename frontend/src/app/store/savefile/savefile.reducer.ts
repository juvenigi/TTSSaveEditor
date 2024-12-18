import {createFeatureSelector, createReducer, on} from '@ngrx/store';

import {initialSaveFileState, SaveFile} from "./savefile.state";
import {SaveFileApiActions} from "./savefile.actions";

export const savefileReducerKey = 'savefileReducer';
export const selectSavefileState = createFeatureSelector<SaveFile>(savefileReducerKey);
export const savefileReducer = createReducer(
  initialSaveFileState,
  on(SaveFileApiActions.requestSavefileByPath, (state, {path}) => {
    return {...state, loadingState: state.path !== path ? 'LOADING' : 'DONE'} satisfies SaveFile
  }),
  on(SaveFileApiActions.savefileSuccess, (_, {data}) => {
    return data;
  }),
);
