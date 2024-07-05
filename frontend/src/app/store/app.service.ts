import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {catchError, map, Observable, take, tap} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class AppService {
  private http = inject(HttpClient);
  private api = "http://localhost:3000/api"

  public getSaveFileJson(path: string): Observable<string> {
    return this.http.get(`${this.api}/savefile`, {params: {path}}).pipe(
      map(o => {
        console.debug(o);
        // console.debug(JSON.stringify(o))
        return JSON.stringify(o)
      })
    )
  }

  public getDirectory(path: string): Observable<string[]> {
    return this.http.get<any>(`${this.api}/directory`, {params: {path}}).pipe(
      map(res => {
        if (!(res != null && typeof res[Symbol.iterator] === 'function')) {
          return [];
        }
        return res.map((entry: any) => entry.path ?? "")
      }),
      catchError(err => {
        console.error(err);
        return []
      }))
  }
}
