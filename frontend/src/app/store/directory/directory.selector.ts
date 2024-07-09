import {selectDirectoryState} from "./directory.reducer";
import {Directory} from "./directory.state";
import {createSelector} from "@ngrx/store";

export const selectRootPath = createSelector(
  selectDirectoryState, (state: Directory) => state.rootPath
);

export const selectRelDirectories = createSelector(
  selectDirectoryState, (state: Directory) => {
    const rootLength = state.rootPath?.length ? state.rootPath.length + 1 : undefined;
    const folders = state.directoryEntries
      .map(entry => entry.slice(rootLength, entry.length))
      .filter(entry => entry.includes(state.relPath))
      .map(entry => entry.slice(state.relPath.length))
      .map(entry => entry.split('\\').filter(i => i)[0])
      .filter(value => value && !value.includes('.'))
      .reduce((accum, value) => accum.add(value), new Set<string>())
    const up = state.relPath.length > 0 ? ['..'] : [];
    console.debug(state.relPath)
    console.debug(state.directoryEntries
      .map(entry => entry.slice(rootLength, entry.length))
      .filter(entry => entry.includes(state.relPath)))
    console.debug(state.directoryEntries
      .map(entry => entry.slice(rootLength, entry.length))
      .filter(entry => entry.includes(state.relPath))
      .map(entry => entry.slice(state.relPath.length)))
    console.debug([...up, ...folders])
    return [...up, ...folders];
  }
);

export const selectPathMask = createSelector(
  selectDirectoryState, (state) => {
    return state.rootPath
      ? [...state.rootPath?.split('\\') ?? '', ...state.relPath.split('\\')].join('\\')
      : "";
  }
);
