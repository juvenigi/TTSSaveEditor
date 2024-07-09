import {createActionGroup, props} from "@ngrx/store";

export const SearchFilterActions = createActionGroup({
  source: 'Search Filter',
  events: {
    'Apply Filter' : props<{ jsonPath?: number[], search?: string}>()
  }
});
