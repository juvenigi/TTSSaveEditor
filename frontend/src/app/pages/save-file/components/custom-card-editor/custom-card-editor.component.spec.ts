import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CustomCardEditorComponent } from './custom-card-editor.component';

describe('CustomCardEditorComponent', () => {
  let component: CustomCardEditorComponent;
  let fixture: ComponentFixture<CustomCardEditorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CustomCardEditorComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CustomCardEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
