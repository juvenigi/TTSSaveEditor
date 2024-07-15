import {createFeatureSelector, createReducer, on} from "@ngrx/store";
import {DirectoryApiActions, IOActions} from "./directory.actions";
import {deduceDividerStyle, Directory, initialDirectoryState} from "./directory.state";

export const directoryReducerKey = 'directoryReducer';
export const selectDirectoryState = createFeatureSelector<Directory>(directoryReducerKey);
export const directoryReducer = createReducer(
  initialDirectoryState,
  on(DirectoryApiActions.requestDirectory, (state: Directory, {rootPath}) => {
    return {...state, loadingState: state.rootPath !== rootPath ? 'LOADING' : state.loadingState} satisfies Directory;
  }),
  on(DirectoryApiActions.directorySuccessResponse, (_: Directory, {directory}) => {
    console.debug(directory);
    return directory;
  }),
  on(DirectoryApiActions.navigateToFolder, (state, {folder}) => {
    const div = deduceDividerStyle(state);
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
    return {...state, relPath: segments.join('\\')} satisfies Directory;
  }),
  on(IOActions.writeSavefileSuccess, (state: Directory, {file}) => {
    const exists = state.directoryEntries.findIndex(value => value === file.path) !== -1;
    return {
      ...state,
      directoryEntries: exists ? state.directoryEntries : [...state.directoryEntries, file.path]
    } satisfies Directory;
  })
)
