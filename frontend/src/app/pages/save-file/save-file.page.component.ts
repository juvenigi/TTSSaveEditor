import {Component, inject} from '@angular/core';
import {Store} from "@ngrx/store";
import {ActivatedRoute, Router} from "@angular/router";
import {filter, firstValueFrom, map, Observable} from "rxjs";
import {AsyncPipe, NgTemplateOutlet} from "@angular/common";
import {NgLetModule} from "ng-let";
import {selectSaveRaw, selectTabletopSave} from "../../store/savefile/savefile.selector";
import {SaveFile} from "../../store/savefile/savefile.models";
import {SaveState} from "../../types/ttstypes";
import {JSONValue} from "../../misc";

@Component({
  selector: 'app-save-file-page',
  standalone: true,
  imports: [
    AsyncPipe,
    NgLetModule,
    NgTemplateOutlet
  ],
  templateUrl: './save-file.page.component.html',
  styleUrl: './save-file.page.component.css'
})
export class SaveFilePageComponent {

  private readonly store = inject(Store)
  private readonly route = inject(ActivatedRoute).snapshot;
  private readonly router = inject(Router);

  raw$: Observable<JSONValue> = this.store.select(selectSaveRaw).pipe(map((save: SaveFile) => save.saveData));
  saveState$ = this.store.select(selectTabletopSave).pipe(
    filter((saveState): saveState is SaveState => saveState !== undefined)
  );
}
