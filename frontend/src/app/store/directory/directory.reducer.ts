import {createFeatureSelector, createReducer, on} from "@ngrx/store";
import {DirectoryApiActions, IOActions} from "./directory.actions";
import {Directory, initialDirectoryState} from "./directory.state";

export const directoryReducer = createReducer(
  initialDirectoryState,
  on(DirectoryApiActions.requestDirectory, (state: Directory, {rootPath}) => {
    return {...state, loadingState: state.rootPath !== rootPath ? 'LOADING' : state.loadingState} satisfies Directory;
  }),
  on(DirectoryApiActions.directorySuccessResponse, (_: Directory, {directory}) => {
    return directory;
  }),
  on(DirectoryApiActions.navigateToFolder, (state, {folder}) => {
    const div = '\\';
    const upDir = folder === ".."
    let segments = state.relPath.split(div);
    segments = segments.slice(0, upDir ? -1 : undefined);
    if (!upDir) {
      if (segments.length === 1 && segments[0] === '') {
        segments = [folder];
      } else {
        segments.push(folder);
      }
    }
    return {...state, relPath: segments.join('\\')};
  }),
  on(IOActions.writeSavefileSuccess, (state: Directory, {file}) => {
    const exists = state.directoryEntries.findIndex(value => file.path) !== -1;
    return {
      ...state,
      directoryEntries: exists ? state.directoryEntries : [...state.directoryEntries, file.path]
    } satisfies Directory;
  })
)
export const directoryReducerKey = 'directoryReducer';
export const selectDirectoryState = createFeatureSelector<Directory>(directoryReducerKey);
