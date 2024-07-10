import {Component, inject} from '@angular/core';
import {Store} from "@ngrx/store";
import {AsyncPipe} from "@angular/common";
import {Router} from "@angular/router";
import {map, Observable} from "rxjs";
import {SaveFileApiActions} from "../../store/savefile/savefile.actions";
import {Directory} from "../../store/directory/directory.state";
import {selectDirectoryState} from "../../store/directory/directory.reducer";
import {selectPathMask, selectRelDirectories} from "../../store/directory/directory.selector";
import {DirectoryApiActions} from "../../store/directory/directory.actions";

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
  private directoryState$ = this.store.select(selectDirectoryState);
  savefiles$: Observable<string[]> = this.directoryState$.pipe(
    map((dir: Directory) => dir.directoryEntries)
  );

  relFolders$: Observable<string[]> = this.store.select(selectRelDirectories)

  pathMask$: Observable<string> = this.store.select(selectPathMask);

  constructor() {

  }

  onSelect(path: string) {
    this.store.dispatch(SaveFileApiActions.requestSavefileByPath({path}))
    void this.router.navigate(["savefile"]);
  }

  onFolderSelect(folder: string) {
    this.store.dispatch(DirectoryApiActions.navigateToFolder({folder}))
  }

}
