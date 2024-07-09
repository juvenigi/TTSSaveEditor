import { ResolveFn } from '@angular/router';

export const cachedEditsResolver: ResolveFn<boolean> = (route, state) => {
  return true;
};
