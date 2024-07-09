import { TestBed } from '@angular/core/testing';
import { CanActivateFn } from '@angular/router';

import { cachedEditsGuard } from './cached-edits.guard';

describe('cachedEditsGuard', () => {
  const executeGuard: CanActivateFn = (...guardParameters) => 
      TestBed.runInInjectionContext(() => cachedEditsGuard(...guardParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it('should be created', () => {
    expect(executeGuard).toBeTruthy();
  });
});
