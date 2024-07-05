import {Component, inject, OnInit} from '@angular/core';
import {Store} from "@ngrx/store";
import {SaveFileApiActions} from "../../store/app.actions";
import {selectAvailablePaths} from "../../store/app.reducer";
import {AsyncPipe} from "@angular/common";
import {async, firstValueFrom, tap} from "rxjs";

@Component({
  selector: 'app-directory-page',
  standalone: true,
  imports: [
    AsyncPipe
  ],
  templateUrl: './directory.page.component.html',
  styleUrl: './directory.page.component.css'
})
export class DirectoryPageComponent implements OnInit{
  store = inject(Store)

  savefiles$ = this.store.select(selectAvailablePaths).pipe(tap((_)=> console.debug("piped")));

  savefile2 = [] as string[]

  ngOnInit() {
    //fixme this is obviously proof-of-concept
    this.store.dispatch(SaveFileApiActions.requestSavefiles({ fsPath: ""}))

    firstValueFrom(this.store.select(selectAvailablePaths)).then(update => this.savefile2 = update);
  }
}
