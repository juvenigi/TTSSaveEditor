import {Component, inject, OnInit} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {Store} from "@ngrx/store";
import {AsyncPipe, NgComponentOutlet} from "@angular/common";
import {firstValueFrom} from "rxjs";
import {AppActions} from "./store/app.actions";
import {HeaderComponent} from "./components/header/header.component";
import {selectFsPath, selectSaveData} from "./store/app.reducer";
import {DirectoryPageComponent} from "./pages/directory.page/directory.page.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, AsyncPipe, HeaderComponent, NgComponentOutlet, DirectoryPageComponent],
  providers: [],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  title = 'frontend';
  private store = inject(Store);

  data$ = this.store.select(selectSaveData);
  fsPath$ = this.store.select(selectFsPath);

  ngOnInit() {
    this.store.dispatch(AppActions.requestSavefile({fsPath: ""}))
  }

  onFSChange($event: Event) {
    console.debug("on change triggered")
    let value = undefined;
    if ($event.target instanceof HTMLInputElement) value = $event.target.value;
    if (!value) return;
    firstValueFrom(this.fsPath$).then((fs) => {
      if (fs !== value) {
        console.debug(`requesting new file: ${value}`)
        this.store.dispatch(AppActions.requestSavefile({fsPath: value}))
      }
    });
  }
}
