import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {map, Observable} from "rxjs";
import {JSONValue} from "../misc";

@Injectable({
  providedIn: 'root'
})
export class IoService {
  private http = inject(HttpClient);
  private api = "http://localhost:3000/api"

  public getSaveFile(path: string): Observable<JSONValue> {
    return this.http.get<JSONValue>(`${this.api}/savefile`, {params: {path}});
  }

  public getDirectory(path: string): Observable<string[]> {
    return this.http.get<any>(`${this.api}/directory`, {params: {path}}).pipe(
      map(res => {
        if (!(res != null && typeof res[Symbol.iterator] === 'function')) {
          return [];
        }
        return res.map((entry: any) => entry.path.toString() ?? "").filter((str: string) => str.length > 0)
      }));
  }
}
