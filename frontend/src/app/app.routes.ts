import {ActivatedRouteSnapshot, CanActivateFn, RouterStateSnapshot, Routes} from '@angular/router';
import {DirectoryPageComponent} from "./pages/directory/directory.page.component";
import {SaveFilePageComponent} from "./pages/save-file/save-file.page.component";
import {inject} from "@angular/core";
import {Store} from "@ngrx/store";
import {selectSaveLoadingState} from "./store/savefile/savefile.selector";
import {firstValueFrom} from "rxjs";

export const unloadedSaveGuard: CanActivateFn = async (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  return await firstValueFrom(inject(Store).select(selectSaveLoadingState))
    .then((state: "DONE" | "LOADING" | "PENDING" | "FAILED") => !['PENDING', 'FAILED'].includes(state));
}

export const routes: Routes = [
  {path: "directory", pathMatch: "full", component: DirectoryPageComponent},
  {path: 'savefile', pathMatch: "full", component: SaveFilePageComponent, canActivate: [unloadedSaveGuard]},
  {path: "**", pathMatch: "full", redirectTo: "directory"},
];

