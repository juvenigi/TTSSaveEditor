import {Component, effect, inject, input, OnDestroy, OnInit} from '@angular/core';
import {GameCardFormData} from "../../../../store/savefile/savefile.state";
import {FormBuilder, ReactiveFormsModule} from "@angular/forms";
import {Observable} from "rxjs";
import {Store} from "@ngrx/store";
import {selectObjectPathMap} from "../../../../store/savefile/savefile.selector";
import {CustomCardActions} from "../../../../store/savefile/savefile.actions";
import {requiredFilter} from "../../../../utils/rxjs.utils";
import {ObjectState} from "../../../../types/ttstypes";

@Component({
  selector: 'app-custom-card-editor',
  standalone: true,
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './custom-card-editor.component.html',
  styleUrl: './custom-card-editor.component.css'
})
export class CustomCardEditorComponent implements OnInit, OnDestroy {
  form = input.required<GameCardFormData>();
  formGroup!: ReturnType<typeof this.initForm>;
  formPatcher = effect(() => {
    this.formGroup.patchValue(this.form());
    this.formGroup.markAsPristine();
  });

  private builder = inject(FormBuilder).nonNullable;
  private store = inject(Store);
  //TODO: add deck/collection selector
  private objects: Observable<Map<string, ObjectState>> = this.store.select(selectObjectPathMap).pipe(requiredFilter());

  ngOnInit() {
    this.formGroup = this.initForm();
  }

  ngOnDestroy() {
    if(this.formGroup.dirty) {
      void this.save()
    }
    this.formPatcher.destroy();
  }

  private initForm() {
    return this.builder.group(this.form());
  }

  async save() {
    console.debug("saveCard fired")
    console.debug(this.formGroup.getRawValue())
    this.store.dispatch(CustomCardActions.pushChanges({cardData: this.formGroup.getRawValue()}))
  }
}
