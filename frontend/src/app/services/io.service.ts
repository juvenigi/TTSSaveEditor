import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {map, Observable} from "rxjs";
import {SaveState} from "../types/ttstypes";

@Injectable({
  providedIn: 'root'
})
export class IoService {
  private http = inject(HttpClient);
  private api = "http://localhost:3000/api"

  public getSaveFile(path: string): Observable<SaveState> {
    return this.http.get<unknown>(`${this.api}/savefile`, {params: {path}}).pipe(map(value => {
        // FIXME: validation
        return value as SaveState;
      })
    );
  }

  public getDirectory(path: string) {
    return this.http.get<{ path: string, entries: Array<{ path: string }> }>(`${this.api}/directory`, {params: {path}})
  };
}
