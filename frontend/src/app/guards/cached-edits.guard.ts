import {CanActivateChildFn, CanActivateFn, CanMatchFn} from '@angular/router';

export const cachedEditsGuard: CanActivateFn = (route, state) => {
  return true;
};
