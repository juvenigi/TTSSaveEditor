import {createFeatureSelector, createReducer, on} from "@ngrx/store";
import {SaveFile} from "../savefile/savefile.state";
import {initialSearchFilterState, SeachFilter} from "./search-filter.state";
import {SearchFilterActions} from "./search-filter.actions";

export const searchFilterReducerKey = 'searchFilterReducer';
export const selectSearchFilter = createFeatureSelector<SeachFilter>(searchFilterReducerKey);
export const searchFilterReducer = createReducer(
  initialSearchFilterState,
  on(SearchFilterActions.applyFilter, (state, {search, jsonPath}) => {
    return Object.assign(state, search, jsonPath);
  })
)
