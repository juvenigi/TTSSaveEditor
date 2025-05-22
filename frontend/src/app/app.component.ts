import {Component, inject, OnInit} from '@angular/core';
import {GoGreet, GoListDir} from '../../wailsjs/go/main/App';
import {Store} from '@ngrx/store';
import {allFilesystemState} from '../store/fs/fs.selector';
import {map, of} from 'rxjs';
import {AsyncPipe} from '@angular/common';

@Component({
  selector: 'app-root',
  imports: [
    AsyncPipe
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  title = 'amazing Person!';



  public ngOnInit() {
    GoGreet('World').then(console.log);
    const promise = GoListDir();

    console.debug("purely angular message")
  }

  protected readonly of = of;
}
