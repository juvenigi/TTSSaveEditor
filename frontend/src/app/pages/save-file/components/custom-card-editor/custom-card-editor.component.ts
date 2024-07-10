import {Component, Input} from '@angular/core';
import {GameCardFormControl} from "../../../../store/savefile/savefile.state";
import {ReactiveFormsModule} from "@angular/forms";

@Component({
  selector: 'app-custom-card-editor',
  standalone: true,
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './custom-card-editor.component.html',
  styleUrl: './custom-card-editor.component.css'
})
export class CustomCardEditorComponent {
  @Input('form') form?: GameCardFormControl;
  get cardForm() {
    return this.form!;
  }


}
