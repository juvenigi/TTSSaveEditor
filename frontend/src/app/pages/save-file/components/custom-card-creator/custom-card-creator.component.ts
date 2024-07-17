import {Component, inject, OnDestroy, OnInit, output} from '@angular/core';
import {initNewCard} from "../../../../store/savefile/savefile.state";
import {FormBuilder, ReactiveFormsModule} from "@angular/forms";
import {map, Unsubscribable, withLatestFrom} from "rxjs";
import {Store} from "@ngrx/store";
import {CustomCardActions} from "../../../../store/savefile/savefile.actions";
import {requiredFilter} from "../../../../utils/rxjs.utils";
import {AsyncPipe, NgClass, NgStyle} from "@angular/common";
import {NgbCollapse} from "@ng-bootstrap/ng-bootstrap";
import {selectSearchFilter} from "../../../../store/savefile-search-filters/search-filter.reducer";
import {Actions, ofType} from "@ngrx/effects";
import {selectAllObjectsPathMap} from "../../../../store/savefile/savefile.selector";

@Component({
  selector: 'app-custom-card-creator',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    NgStyle,
    AsyncPipe,
    NgbCollapse,
    NgClass
  ],
  templateUrl: './custom-card-creator.component.html',
  styleUrl: './custom-card-creator.component.css'
})
export class CustomCardCreatorComponent implements OnInit, OnDestroy {
  submit = output<boolean>();

  formattedCardPath = "";
  formGroup!: ReturnType<typeof this.initForm>;

  submitting = false;
  addingToDeck = false;

  private me = Symbol("reference to this component");

  private builder = inject(FormBuilder).nonNullable;
  private store = inject(Store);
  private actions = inject(Actions);

  private subscriptions = [] as Unsubscribable[];
  private pathFilter$ = this.store.select(selectSearchFilter)
    .pipe(requiredFilter(), map(({jsonPath}) => jsonPath.join('/')));

  ngOnInit() {
    this.formGroup = this.initForm('');
    this.subscriptions.push(
      this.pathFilter$.pipe(withLatestFrom(this.store.select(selectAllObjectsPathMap)))
        .subscribe(([pathUpdate, objects]) => {
          console.debug('debugging', pathUpdate, objects?.get(pathUpdate))
          // adding to deck validation
          this.addingToDeck = (objects?.has(pathUpdate) ?? false) && objects!.get(pathUpdate)!.Name === 'Deck';

          // path
          let parents = [] as string[];
          const splitCardPath = pathUpdate.split('/');
          for (let i = 1; i <= splitCardPath.length; i++) {
            parents.push(splitCardPath.slice(0, i).join('/'));
          }
          this.formGroup.controls['path'].setValue(pathUpdate);
          console.debug(parents);
          this.formattedCardPath = parents
            .map(p => objects?.get(p)?.Nickname ?? 'unknown')
            .join(" / ").toUpperCase()
        }),

      this.actions.pipe(
        ofType(CustomCardActions.pushNewSuccess, CustomCardActions.pushNewFailure),
        map(({ref}) => ref)
      ).subscribe(ref => {
        if (ref === this.me) {
          this.resetToBlank()
          this.submit.emit(true);
        }
      }),
    );
  }

  ngOnDestroy() {
    console.debug("component destroyed: " + this.formGroup.value.path ?? 'unknown path')
    if (this.formGroup.dirty) {
      void this.save()
    }
    this.subscriptions.forEach(i => i.unsubscribe());
  }

  async save() {
    this.submitting = true;
    this.store.dispatch(CustomCardActions.pushNewCard({cardData: this.formGroup.getRawValue(), ref: this.me}));
  }

  private initForm(path: string) {
    return this.builder.group(initNewCard(path));
  }

  private resetToBlank() {
    this.formGroup.controls['cardText'].reset();
    this.formGroup.markAsPristine();
  }
}
