import {Component, inject, OnDestroy, OnInit} from '@angular/core';
import {Store} from "@ngrx/store";
import {AsyncPipe, NgClass} from "@angular/common";
import {debounceTime, distinctUntilChanged, firstValueFrom, map, Unsubscribable, withLatestFrom} from "rxjs";
import {ActivatedRoute, Router} from "@angular/router";
import {NgLetModule} from "ng-let";
import {selectSaveName} from "../../store/savefile/savefile.selector";
import {selectRootPath} from "../../store/directory/directory.selector";
import {SAVEFILE_ROUTE} from "../../app.routes";
import {NgbCollapse} from "@ng-bootstrap/ng-bootstrap";
import {DirectoryApiActions} from "../../store/directory/directory.actions";
import {FormControl, ReactiveFormsModule} from "@angular/forms";

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    AsyncPipe,
    NgClass,
    NgLetModule,
    NgbCollapse,
    ReactiveFormsModule
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.css'
})
export class HeaderComponent implements OnDestroy {
  private readonly store = inject(Store)
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private subscriptions: Unsubscribable [] = [];
  saveName$ = this.store.select(selectSaveName);
  rootPath$ = this.store.select(selectRootPath);
  saveFileActive$ = this.route.url.pipe(map(segments =>
    segments.findIndex(segment => segment.path === SAVEFILE_ROUTE) !== -1
  ));

  rootDirForm: FormControl<string> = new FormControl<string>("", {nonNullable: true});
  rootDirSelectorCollapsed = true;

  constructor() {
    void this.setupForm();
  }

  private async setupForm() {
    this.subscriptions.push(
      this.rootPath$.subscribe(update => {
        this.rootDirSelectorCollapsed = ((update ?? '').length > 0);
        this.rootDirForm.patchValue(update ?? '', {emitEvent: false});
        this.rootDirForm.markAsPristine({emitEvent: false})
      }),
      this.rootDirForm.valueChanges
        .pipe(distinctUntilChanged(), debounceTime(3000), withLatestFrom(this.rootPath$))
        .subscribe(([rootPath, oldRootPath]) => {
          console.debug(rootPath);

          if (rootPath.length > 0 && rootPath !== (oldRootPath ?? '')) {
            this.store.dispatch(DirectoryApiActions.requestDirectory({rootPath}))
            void this.router.navigate(['directory']);
          }
          this.rootDirSelectorCollapsed = false;
        })
    );
  }

  ngOnDestroy() {
    this.subscriptions.forEach(i => i.unsubscribe());
  }

  navigateToDirectory() {
    void this.router.navigate(['directory']);
  }
}
