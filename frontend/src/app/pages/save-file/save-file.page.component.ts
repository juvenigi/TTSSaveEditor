import {Component, inject, OnDestroy} from '@angular/core';
import {Store} from "@ngrx/store";
import {debounceTime, distinctUntilChanged, Observable, Unsubscribable} from "rxjs";
import {AsyncPipe, NgTemplateOutlet} from "@angular/common";
import {NgLetModule} from "ng-let";
import {selectCardForms, selectFlattenedObjects, selectSavefileHeader} from "../../store/savefile/savefile.selector";
import {requiredFilter} from "../../utils/rxjs.utils";
import {FormControl, ReactiveFormsModule} from "@angular/forms";
import {SearchFilterActions} from "../../store/search-filters/search-filter.actions";
import {CustomCardEditorComponent} from "./components/custom-card-editor/custom-card-editor.component";
import {tryCardInit} from "../../store/savefile/savefile.state";

@Component({
  selector: 'app-save-file-page',
  standalone: true,
  imports: [
    AsyncPipe,
    NgLetModule,
    NgTemplateOutlet,
    ReactiveFormsModule,
    CustomCardEditorComponent
  ],
  templateUrl: './save-file.page.component.html',
  styleUrl: './save-file.page.component.css'
})
export class SaveFilePageComponent implements OnDestroy {
  private readonly store = inject(Store)

  saveState$ = this.store.select(selectSavefileHeader).pipe(requiredFilter());

  private flat$ = this.store.select(selectFlattenedObjects).pipe(requiredFilter());
  flatObjects$ = this.flat$;
  cardForms$: Observable<Map<string, Exclude<ReturnType<typeof tryCardInit>, undefined>>> = this.store.select(selectCardForms).pipe(requiredFilter())

  objectSearch = new FormControl<string>('', {nonNullable: true});
  private subscriptions = [] as Unsubscribable[];

  constructor() {
    this.subscriptions.push(
      this.objectSearch.valueChanges
        .pipe(distinctUntilChanged(), debounceTime(200))
        .subscribe(
          search => this.store.dispatch(SearchFilterActions.applyFilter({search}))
        )
    );
  }

  ngOnDestroy() {
    this.subscriptions.forEach(i => i.unsubscribe());
  }
}
