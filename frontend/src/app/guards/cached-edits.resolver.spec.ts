import { TestBed } from '@angular/core/testing';
import { ResolveFn } from '@angular/router';

import { cachedEditsResolver } from './cached-edits.resolver';

describe('cachedEditsResolver', () => {
  const executeResolver: ResolveFn<boolean> = (...resolverParameters) => 
      TestBed.runInInjectionContext(() => cachedEditsResolver(...resolverParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it('should be created', () => {
    expect(executeResolver).toBeTruthy();
  });
});
