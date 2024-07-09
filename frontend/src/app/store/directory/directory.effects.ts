import {Actions, createEffect, ofType} from "@ngrx/effects";
import {DirectoryApiActions} from "./directory.actions";
import {catchError, map, of, switchMap} from "rxjs";
import {inject} from "@angular/core";
import {IoService} from "../../services/io.service";

export class DirectoryEffects {
  actions$ = inject(Actions)
  service = inject(IoService)

  fetchDir$ = createEffect(() => this.actions$.pipe(
    ofType(DirectoryApiActions.requestDirectory),
    switchMap(({rootPath}) => this.service.getDirectory(rootPath).pipe(
      map((directoryEntries: string[]) => DirectoryApiActions
        .directorySuccessResponse({directory: {directoryEntries, rootPath, relPath: "", loadingState: "DONE"}})),
      catchError((err) => of(DirectoryApiActions.directoryFailureResponse({request: rootPath, message: err})))
    ))
  ));
}

