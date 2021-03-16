import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectMachinesComponent } from './select-machines.component';

describe('SelectMachinesComponent', () => {
  let component: SelectMachinesComponent;
  let fixture: ComponentFixture<SelectMachinesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectMachinesComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectMachinesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
