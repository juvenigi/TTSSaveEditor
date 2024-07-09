import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SaveFilePageComponent } from './save-file.page.component';

describe('SaveFilePageComponent', () => {
  let component: SaveFilePageComponent;
  let fixture: ComponentFixture<SaveFilePageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SaveFilePageComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SaveFilePageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
