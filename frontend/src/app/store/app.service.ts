import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {map, Observable, take, tap} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class AppService {
  private http = inject(HttpClient);

  public getSaveFileJson(path: string): Observable<string> {
    return this.http.get("http://localhost:3000/api/savefile", {params: {path}}).pipe(
      map( o => {
        console.debug(o);
        // console.debug(JSON.stringify(o))
        return JSON.stringify(o)
      })
    )
  }
}
