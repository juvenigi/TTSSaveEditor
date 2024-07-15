import {inject, Injectable} from "@angular/core";
import {Actions, createEffect, ofType} from "@ngrx/effects";
import {IoService} from "../../services/io.service";
import {catchError, EMPTY, map, of, switchMap, withLatestFrom} from "rxjs";
import {Store} from "@ngrx/store";
import {CustomCardActions, SaveFileApiActions} from "./savefile.actions";
import {getObjectPatch, SaveFile} from "./savefile.state";
import {SaveState} from "../../types/ttstypes";
import {selectObjectPathMap, selectSavefileHeader} from "./savefile.selector";
import {SearchFilterActions} from "../savefile-search-filters/search-filter.actions";

@Injectable()
export class SavefileEffects {
  actions$ = inject(Actions)
  store: Store = inject(Store)
  io = inject(IoService)


  resetFilterOnLoad$ = createEffect(() => this.actions$.pipe(
    ofType(SaveFileApiActions.requestSavefileByPath),
    withLatestFrom(this.store.select(selectSavefileHeader)),
    switchMap(([{path}, old]) => {
      const oldPath = old?.path ?? "";
      const updatePath = path;
      console.debug(oldPath, updatePath);
      return oldPath === updatePath ? EMPTY : of(SearchFilterActions.applyFilter({jsonPath: [], search: ""}))
    })
  ));

  fetchFile$ = createEffect(() => this.actions$.pipe(
    ofType(SaveFileApiActions.requestSavefileByPath),
    switchMap(({path}) => {
      return this.io
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

  pushPatches$ = createEffect(() => this.actions$.pipe(
    ofType(CustomCardActions.pushChanges),
    withLatestFrom(
      this.store.select(selectObjectPathMap),
      this.store.select(selectSavefileHeader)
    ),
    switchMap(
      ([
         {cardData},
         objMap,
         header
       ]) => {
        const original = objMap?.get(cardData.path)
        const savePath = header?.path;
        if (!original || !savePath) return of(CustomCardActions.pushChangesFail({reason: "invalid savestate"}));

        const patch = getObjectPatch(cardData, original);
        if (!patch) return of(CustomCardActions.pushChangesFail({reason: "could not create patch"}));

        return this.io.pushPatches(savePath, patch).pipe(
          map((data: SaveState) => SaveFileApiActions.savefileSuccess({
            data: ({
              saveData: data,
              path: savePath,
              objectEdits: new Map(),
              loadingState: "DONE"
            } satisfies SaveFile)
          })),
          catchError(() => of(CustomCardActions.pushChangesFail({reason: "network error"})))
        );
      })
  ));
}

