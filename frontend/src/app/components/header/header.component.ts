import {Component, inject, OnDestroy} from '@angular/core';
import {Store} from "@ngrx/store";
import {AsyncPipe, NgClass} from "@angular/common";
import {debounceTime, distinctUntilChanged, map, Unsubscribable, withLatestFrom} from "rxjs";
import {NavigationEnd, Router, RouterLink, RouterLinkActive} from "@angular/router";
import {NgLetModule} from "ng-let";
import {selectSaveName} from "../../store/savefile/savefile.selector";
import {selectRootPath} from "../../store/directory/directory.selector";
import {NgbCollapse} from "@ng-bootstrap/ng-bootstrap";
import {DirectoryApiActions} from "../../store/directory/directory.actions";
import {FormControl, ReactiveFormsModule} from "@angular/forms";
import {instanceOfFilter} from "../../utils/rxjs.utils";

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    AsyncPipe,
    NgClass,
    NgLetModule,
    NgbCollapse,
    ReactiveFormsModule,
    RouterLink,
    RouterLinkActive
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.css'
})
export class HeaderComponent implements OnDestroy {
  private readonly store = inject(Store)
  private readonly router = inject(Router);
  private subscriptions: Unsubscribable [] = [];
  saveName$ = this.store.select(selectSaveName);
  rootPath$ = this.store.select(selectRootPath);
  currentPath$ = this.router.events.pipe(
    instanceOfFilter(NavigationEnd),
    map((event) => event.urlAfterRedirects)
  );


  rootDirForm: FormControl<string> = new FormControl<string>("", {nonNullable: true});
  rootDirSelectorCollapsed = true;

  constructor() {
    this.store.dispatch(DirectoryApiActions.requestDirectory({rootPath: this.rootDirForm.value}));
    this.setupForm();
    this.setupFormToggle();
  }

  /**
   * Hide input once we are no longer in the '/directory' route
   * @private
   */
  private setupFormToggle() {
    this.subscriptions.push(
      this.currentPath$.subscribe(path => {
        if (path !== "/directory") {
          this.rootDirSelectorCollapsed = true;
        }
      })
    )
  }

  private setupForm() {
    this.subscriptions.push(
      this.rootPath$.subscribe(update => {
        this.rootDirSelectorCollapsed = (update ?? '').length === 0 ? false : this.rootDirSelectorCollapsed;
        this.rootDirForm.patchValue(update ?? '', {emitEvent: false});
        this.rootDirForm.markAsPristine({emitEvent: false})
      }),
      this.rootDirForm.valueChanges
        .pipe(distinctUntilChanged(), debounceTime(400), withLatestFrom(this.rootPath$))
        .subscribe(([rootPath, oldRootPath]) => {
          console.debug(rootPath, oldRootPath);
          if (rootPath.length > 0 && rootPath !== (oldRootPath ?? '')) {
            this.store.dispatch(DirectoryApiActions.requestDirectory({rootPath}))
            void this.router.navigate(['directory']);
          }
        })
    );
  }

  ngOnDestroy() {
    this.subscriptions.forEach(i => i.unsubscribe());
  }

  onSubmit(event: Event) {
    this.rootDirSelectorCollapsed = true;
  }
}
