import {createFeatureSelector, createReducer, on} from "@ngrx/store";
import {DirectoryApiActions} from "./directory.actions";
import {Directory, initialDirectoryState} from "./directory.models";

export const directoryReducer = createReducer(
  initialDirectoryState,
  on(DirectoryApiActions.requestDirectory, (state: Directory, {rootPath}) => {
    return {...state, loadingState: state.rootPath !== rootPath ? 'LOADING' : state.loadingState} satisfies Directory;
  }),
  on(DirectoryApiActions.directorySuccessResponse, (_: Directory, {directory}) => {
    return directory;
  }),
  on(DirectoryApiActions.writeSavefileSuccess, (state:Directory, {file})=>{
    const exists = state.directoryEntries.findIndex(value => file.path) !== -1;
    return {...state, directoryEntries: exists? state.directoryEntries : [...state.directoryEntries, file.path]} satisfies Directory;
  })
)
export const directoryReducerKey = 'directoryReducer';
export const selectDirectoryState = createFeatureSelector<Directory>(directoryReducerKey);
