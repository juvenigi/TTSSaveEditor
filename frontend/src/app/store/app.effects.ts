import {inject, Injectable} from "@angular/core";
import {Actions, createEffect, ofType} from "@ngrx/effects";
import {AppService} from "./app.service";
import {catchError, EMPTY, exhaustMap, map, switchMap} from "rxjs";
import {AppActions, SaveFileApiActions} from "./app.actions";

@Injectable()
export class AppEffects {
  actions$ = inject(Actions)
  service = inject(AppService)

  fetchFile$ = createEffect(() => this.actions$.pipe(
    ofType(AppActions.requestSavefile),
    switchMap(({fsPath}) => {
      console.debug("effect fired")
      return this.service
        .getSaveFileJson(fsPath).pipe(
          map((data: string) => {
            console.debug("fresh data")
            return {type: '[SaveFile API] Savefile Retrieve Success', data, fsPath};
          }),
          catchError(() => EMPTY)
        );
    })
  ));

  fetchDir$ = createEffect(() => this.actions$.pipe(
    ofType(SaveFileApiActions.requestSavefiles),
    switchMap(({fsPath}) => this.service.getDirectory(fsPath).pipe(
      map((paths: string[]) => {
        console.debug("effect fired")
        console.log(fsPath)
        console.debug(paths)
        return {type: '[SaveFile API] Request Savefiles Success', paths, fsPath }
      }),
      catchError(() => EMPTY)
    ))
  ));
}
