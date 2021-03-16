import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CustomSwitchComponent } from './customswitch.component';

describe('CustomSwitchComponent', () => {
  let component: CustomSwitchComponent;
  let fixture: ComponentFixture<CustomSwitchComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CustomSwitchComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CustomSwitchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
