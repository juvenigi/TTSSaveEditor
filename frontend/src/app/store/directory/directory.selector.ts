import {selectDirectoryState} from "./directory.reducer";
import {Directory} from "./directory.models";
import {createSelector} from "@ngrx/store";

export const selectDirectory = selectDirectoryState;

export const selectRootPath = createSelector(
  selectDirectoryState, (state: Directory) => state.rootPath
);
