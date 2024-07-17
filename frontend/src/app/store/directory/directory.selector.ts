import {selectDirectoryState} from "./directory.reducer";
import {deduceDividerStyle, Directory} from "./directory.state";
import {createSelector} from "@ngrx/store";

export const selectRootPath = createSelector(
  selectDirectoryState, (state: Directory) => state.rootPath
);

export const selectRelDirectories = createSelector(
  selectDirectoryState, (state: Directory) => {
    const div = deduceDividerStyle(state);
    const rootLength = state.rootPath?.length ? state.rootPath.length + 1 : undefined;
    const folders = state.directoryEntries
      .map(entry => entry.slice(rootLength, entry.length))
      .filter(entry => entry.includes(state.relPath))
      .map(entry => entry.slice(state.relPath.length))
      .map(entry => entry.split(div).filter(i => i)[0])
      .filter(value => value && !value.includes('.'))
      .reduce((accum, value) => accum.add(value), new Set<string>())
    const up = state.relPath.length > 0 ? ['..'] : [];
    return [...up, ...[...folders].sort()];
  }
);

export const selectPathMask = createSelector(
  selectDirectoryState, (state) => {
    const div = deduceDividerStyle(state);
    return state.rootPath
      ? [...state.rootPath?.split(div) ?? '', ...state.relPath.split(div)].join(div)
      : "";
  }
);
