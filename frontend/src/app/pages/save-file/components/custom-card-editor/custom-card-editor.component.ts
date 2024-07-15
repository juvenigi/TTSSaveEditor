import {Component, inject, Input, OnInit} from '@angular/core';
import {GameCardFormData} from "../../../../store/savefile/savefile.state";
import {FormBuilder, ReactiveFormsModule} from "@angular/forms";
import CardService from "../../../../services/card.service";
import {firstValueFrom, tap} from "rxjs";
import {Store} from "@ngrx/store";
import {selectObjectPathMap} from "../../../../store/savefile/savefile.selector";

@Component({
  selector: 'app-custom-card-editor',
  standalone: true,
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './custom-card-editor.component.html',
  styleUrl: './custom-card-editor.component.css'
})
export class CustomCardEditorComponent implements OnInit {
  @Input('form') formData!: GameCardFormData;
  private builder = inject(FormBuilder).nonNullable;
  private cardService = inject(CardService);
  private store = inject(Store);
  objects = this.store.select(selectObjectPathMap).pipe(tap(console.debug));
  form!: ReturnType<typeof this.initForm>

  ngOnInit() {
    this.form = this.builder.group(this.formData)
  }

  private initForm() {
    return this.builder.group(this.formData);
  }

  get cardForm() {
    return this.form.controls.cardScript;
  }

  get scriptForm() {
    return this.form.controls.cardText;
  }

  async saveCard() {
    console.debug("saveCard fired")
    const objPath = this.form.value.path;
    if(!objPath || !this.form.value) return;
    const original = await firstValueFrom(this.objects)
      .then(item => item?.get(objPath))
    if (!original) return;
    console.debug(original)
    console.debug(this.cardService.createEditPatch(this.form.getRawValue(), original));
  }
}
