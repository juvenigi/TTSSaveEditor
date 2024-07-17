import {ComponentFixture, TestBed} from '@angular/core/testing';

import {CustomCardCreatorComponent} from './custom-card-creator.component';

describe('CustomCardEditorComponent', () => {
  let component: CustomCardCreatorComponent;
  let fixture: ComponentFixture<CustomCardCreatorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CustomCardCreatorComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CustomCardCreatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
