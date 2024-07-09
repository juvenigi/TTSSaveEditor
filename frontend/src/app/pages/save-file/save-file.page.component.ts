import {Component, inject} from '@angular/core';
import {Store} from "@ngrx/store";
import {ActivatedRoute, Router} from "@angular/router";
import {map} from "rxjs";
import {AsyncPipe, NgTemplateOutlet} from "@angular/common";
import {NgLetModule} from "ng-let";
import {selectFlattenedObjects, selectSavefileHeader} from "../../store/savefile/savefile.selector";
import {requiredFilter} from "../../utils/rxjs.utils";

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

  saveState$ = this.store.select(selectSavefileHeader).pipe(requiredFilter());

  private flat$ = this.store.select(selectFlattenedObjects).pipe(requiredFilter());
  flatObjects$ = this.flat$.pipe(map((objs) => objs.map(o => o.object)));
}
