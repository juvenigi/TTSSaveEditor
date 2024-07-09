import {Component, inject, OnInit} from '@angular/core';
import {Store} from "@ngrx/store";
import {AsyncPipe} from "@angular/common";
import {Router} from "@angular/router";
import {selectDirectory} from "../../store/directory/directory.selector";
import {map, Observable, tap} from "rxjs";
import {SaveFileApiActions} from "../../store/savefile/savefile.actions";
import {Directory} from "../../store/directory/directory.models";

@Component({
  selector: 'app-directory-page',
  standalone: true,
  imports: [
    AsyncPipe
  ],
  templateUrl: './directory.page.component.html',
  styleUrl: './directory.page.component.css'
})
export class DirectoryPageComponent {
  private readonly store = inject(Store);
  private readonly router = inject(Router);

  savefiles$: Observable<string[]> = this.store.select(selectDirectory)
    .pipe(map((dir: Directory) => dir.directoryEntries));

  onSelect(path: string) {
    this.store.dispatch(SaveFileApiActions.requestSavefileByPath({path}))
    void this.router.navigate(["savefile"]);
  }
}
