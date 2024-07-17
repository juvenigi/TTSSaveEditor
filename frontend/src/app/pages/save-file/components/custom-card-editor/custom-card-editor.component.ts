import {Component, inject, input, OnDestroy, OnInit} from '@angular/core';
import {GameCardFormData} from "../../../../store/savefile/savefile.state";
import {FormBuilder, ReactiveFormsModule} from "@angular/forms";
import {firstValueFrom, map, Observable, Unsubscribable} from "rxjs";
import {Store} from "@ngrx/store";
import {selectAllObjectsPathMap, selectCardForms} from "../../../../store/savefile/savefile.selector";
import {CustomCardActions} from "../../../../store/savefile/savefile.actions";
import {requiredFilter} from "../../../../utils/rxjs.utils";
import {ObjectState} from "../../../../types/ttstypes";
import {AsyncPipe, NgClass, NgStyle} from "@angular/common";
import {NgbCollapse} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-custom-card-editor',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    NgStyle,
    AsyncPipe,
    NgbCollapse,
    NgClass
  ],
  templateUrl: './custom-card-editor.component.html',
  styleUrl: './custom-card-editor.component.css'
})
export class CustomCardEditorComponent implements OnInit, OnDestroy {
  cardPath = input.required<string>();
  formattedCardPath = "";
  // FIXME: I could have used the change if initialValue, but the parent template is conservative about firing the change detection
  // and I don't know how to fix it without causing the component to be destroyed and recreated
  initialValue = input.required<GameCardFormData>();
  formGroup!: ReturnType<typeof this.initForm>;

  submitting = false;
  collapsed = {
    script: true,
    text: false
  }

  private builder = inject(FormBuilder).nonNullable;
  private store = inject(Store);
  private cardForms$: Observable<Map<string, GameCardFormData>> = this.store.select(selectCardForms).pipe(requiredFilter());
  private currentCard$ = this.cardForms$
    .pipe(map(forms => forms.get(this.cardPath())), requiredFilter());
  private objects: Observable<Map<string, ObjectState>> = this.store.select(selectAllObjectsPathMap).pipe(requiredFilter());
  private subscriptions = [] as Unsubscribable[];

  deleteHidden = true;

  ngOnInit() {
    this.formGroup = this.initForm();
    this.subscriptions.push(
      this.currentCard$.subscribe(formValue => {
        if (this.formGroup.pristine || this.submitting) {
          this.submitting = false;
          this.formGroup.patchValue(formValue);
          this.formGroup.markAsPristine();
        }
      })
    );
    Promise.all([firstValueFrom(this.currentCard$), firstValueFrom(this.objects)])
      .then(([card, objects]) => {
        let parents = [] as string[];
        const splitCardPath = card.path.split('/');
        for (let i = 1; i < splitCardPath.length; i++) {
          parents.push(splitCardPath.slice(0, i).join('/'));
        }
        this.formattedCardPath = parents.map(p => objects.get(p)?.Nickname ?? 'unknown').join(" / ").toUpperCase();
      });
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
    this.store.dispatch(CustomCardActions.pushChanges({cardData: this.formGroup.getRawValue()}));
  }

  async revert() {
    const value = await firstValueFrom(this.currentCard$);
    this.formGroup.patchValue(value);
    this.formGroup.markAsPristine();
  }

  selectTab(section: keyof typeof this.collapsed) {
    Object.keys(this.collapsed).forEach((key) => (this.collapsed as any)[key] = true);
    this.collapsed[section] = false;
  }

  private initForm() {
    const group = this.builder.group(this.initialValue());
    group.controls['cardScript'].disable();
    return group;
  }

  deleteMe() {
    this.store.dispatch(CustomCardActions.deleteCard({itemPath: this.formGroup.getRawValue().path}))
  }
}
