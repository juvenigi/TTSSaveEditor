import {ApplicationConfig, provideZoneChangeDetection} from '@angular/core';
import {provideRouter} from '@angular/router';

import {routes} from './app.routes';
import {provideStore, StoreModule} from '@ngrx/store';
import {provideEffects} from '@ngrx/effects';
import {AppEffects} from "./store/app.effects";
import {provideHttpClient} from "@angular/common/http";
import {reducer, reducerFeatureKey} from "./store/app.reducer";

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({eventCoalescing: true}),
    provideRouter(routes),
    provideStore({[reducerFeatureKey]: reducer}),
    provideEffects(AppEffects),
    provideHttpClient(),
  ],
};
