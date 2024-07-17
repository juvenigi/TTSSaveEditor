import {Component, inject, OnDestroy} from '@angular/core';
import {Store} from "@ngrx/store";
import {debounceTime, distinctUntilChanged, firstValueFrom, map, Observable, tap, Unsubscribable} from "rxjs";
import {AsyncPipe, KeyValuePipe, NgClass, NgTemplateOutlet} from "@angular/common";
import {NgLetModule} from "ng-let";
import {
  selectCardForms,
  selectCards,
  selectCollections,
  selectSavefileHeader
} from "../../store/savefile/savefile.selector";
import {requiredFilter} from "../../utils/rxjs.utils";
import {FormControl, ReactiveFormsModule} from "@angular/forms";
import {SearchFilterActions} from "../../store/savefile-search-filters/search-filter.actions";
import {CustomCardEditorComponent} from "./components/custom-card-editor/custom-card-editor.component";
import {GameCardFormData} from "../../store/savefile/savefile.state";
import {selectSearchFilter} from "../../store/savefile-search-filters/search-filter.reducer";
import {NgbCollapse} from "@ng-bootstrap/ng-bootstrap";
import {CustomCardCreatorComponent} from "./components/custom-card-creator/custom-card-creator.component";
import {AutoFocusDirective} from "../../directives/auto-focus-directive";

@Component({
  selector: 'app-save-file-page',
  standalone: true,
  imports: [
    AsyncPipe,
    NgLetModule,
    NgTemplateOutlet,
    ReactiveFormsModule,
    CustomCardEditorComponent,
    KeyValuePipe,
    NgbCollapse,
    CustomCardCreatorComponent,
    NgClass,
    AutoFocusDirective
  ],
  templateUrl: './save-file.page.component.html',
  styleUrl: './save-file.page.component.css'
})
export class SaveFilePageComponent implements OnDestroy {
  private readonly store = inject(Store)

  saveState$ = this.store.select(selectSavefileHeader).pipe(requiredFilter());
  collections$ = this.store.select(selectCollections).pipe(requiredFilter());
  cards$ = this.store.select(selectCards).pipe(requiredFilter(), tap(cards => console.debug(cards)));

  pathFilter$ = this.store.select(selectSearchFilter)
    .pipe(requiredFilter(), map(({jsonPath}) => jsonPath));
  cardForms$: Observable<Map<string, GameCardFormData>> = this.store.select(selectCardForms)
    .pipe(requiredFilter(), distinctUntilChanged())

  objectSearch = new FormControl<string>('', {nonNullable: true});
  private subscriptions = [] as Unsubscribable[];
  collapseCardCreator = true;

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

  collectionPressed(jsonPath: number[]) {
    this.store.dispatch(SearchFilterActions.applyFilter({jsonPath}));
  }

  goUpCollection() {
    firstValueFrom(this.pathFilter$)
      .then((path: number[]) => this.store.dispatch(
        SearchFilterActions.applyFilter({jsonPath: path.slice(0, path.length - 1)})
      ));
  }

  handleSubmitState($event: boolean) {
    if ($event) {
      this.collapseCardCreator = true;
    }
  }
}
