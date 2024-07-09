import {inject, Injectable} from "@angular/core";
import {Actions, createEffect, ofType} from "@ngrx/effects";
import {IoService} from "../../services/io.service";
import {catchError, EMPTY, map, switchMap} from "rxjs";
import {Store} from "@ngrx/store";
import {SaveFileApiActions} from "./savefile.actions";
import {SaveFile} from "./savefile.state";
import {JSONValue} from "../../types/json";
import {SaveState} from "../../types/ttstypes";

@Injectable()
export class SavefileEffects {
  actions$ = inject(Actions)
  store: Store = inject(Store)
  service = inject(IoService)

  fetchFile$ = createEffect(() => this.actions$.pipe(
    ofType(SaveFileApiActions.requestSavefileByPath),
    switchMap(({path}) => {
      return this.service
        .getSaveFile(path).pipe(
          map((data: SaveState) => SaveFileApiActions.savefileSuccess({
            data: ({
              saveData: data,
              path,
              objectEdits: new Map(),
              loadingState: "DONE"
            } satisfies SaveFile)
          })),
          catchError(() => EMPTY)
        );
    })
  ));

}
