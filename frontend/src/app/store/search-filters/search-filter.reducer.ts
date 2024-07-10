import {createFeatureSelector, createReducer, on} from "@ngrx/store";
import {SaveFile} from "../savefile/savefile.state";
import {initialSearchFilterState, SearchFilter} from "./search-filter.state";
import {SearchFilterActions} from "./search-filter.actions";

export const searchFilterReducerKey = 'searchFilterReducer';
export const selectSearchFilter = createFeatureSelector<SearchFilter>(searchFilterReducerKey);
export const searchFilterReducer = createReducer(
  initialSearchFilterState,
  on(SearchFilterActions.applyFilter, (state: SearchFilter, {search, jsonPath}) => {
    return {...state, search: search ?? state.search, jsonPath: jsonPath ?? state.jsonPath} satisfies SearchFilter
  })
)
