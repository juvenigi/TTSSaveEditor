import {Component, Input} from '@angular/core';
import {tryCardInit} from "../../../../store/savefile/savefile.state";
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
  @Input('form') form: ReturnType<typeof tryCardInit>;
  get cardForm() {
    return this.form!;
  }


}
