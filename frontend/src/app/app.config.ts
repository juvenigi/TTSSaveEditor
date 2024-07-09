import {ApplicationConfig, provideZoneChangeDetection} from '@angular/core';
import {provideRouter} from '@angular/router';

import {routes} from './app.routes';
import {provideStore, StoreModule} from '@ngrx/store';
import {provideEffects} from '@ngrx/effects';
import {SavefileEffects} from "./store/savefile/savefile.effects";
import {provideHttpClient} from "@angular/common/http";
import {savefileReducer, savefileReducerKey} from "./store/savefile/savefile.reducer";
import {provideRouterStore} from '@ngrx/router-store';
import {directoryReducer, directoryReducerKey} from "./store/directory/directory.reducer";
import {DirectoryEffects} from "./store/directory/directory.effects";
import {searchFilterReducer, searchFilterReducerKey} from "./store/search-filters/search-filter.reducer";

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({eventCoalescing: true}),
    provideRouter(routes),
    provideStore({
      [savefileReducerKey]: savefileReducer,
      [directoryReducerKey]: directoryReducer,
      [searchFilterReducerKey]: searchFilterReducer
    }),
    provideEffects(SavefileEffects, DirectoryEffects),
    provideHttpClient(),
    provideRouterStore()
  ],
};
