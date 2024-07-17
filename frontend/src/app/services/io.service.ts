import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {map, Observable} from "rxjs";
import {ObjectState, SaveState} from "../types/ttstypes";
import {Patch} from "generate-json-patch";

@Injectable({
  providedIn: 'root'
})
export class IoService {
  private http = inject(HttpClient);
  private api = "http://localhost:3000/api"

// FIXME: response validation

  public getSaveFile(path: string): Observable<SaveState> {
    return this.http.get<unknown>(`${this.api}/savefile`, {params: {path}}).pipe(map(value => {
        return value as SaveState;
      })
    );
  }

  public getDirectory(path: string) {
    return this.http.get<{ path: string, entries: Array<{ path: string }> }>(`${this.api}/directory`, {params: {path}})
  };

  public pushPatches(path: string, patch: Patch) {
    return this.http.patch<unknown>(`${this.api}/savefile`, patch, {params: {path}}).pipe(map(value => {
      return value as SaveState;
    }));
  }

  pushNewCard(savepath: string, jsonPath: string, LuaScript: string, LuaScriptState: string) {
    return this.http.post(`${this.api}/card`, {LuaScript, LuaScriptState} satisfies Partial<ObjectState>,
      {params: {'path': savepath, jsonPath}})
      .pipe(map(value => value as SaveState));
  }

  deleteCard(path: string, jsonPath: string) {
    return this.http.delete(`${this.api}/card`, {params: {path, jsonPath}}).pipe(map(value => value as SaveState));
  }
}
