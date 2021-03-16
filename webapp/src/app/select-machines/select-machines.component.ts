import { Component, OnInit,  ChangeDetectionStrategy } from '@angular/core';
import { NG_VALUE_ACCESSOR, ControlValueAccessor } from '@angular/forms';
import { DataService } from '../data.service';
import { Observable } from 'rxjs';
import { Device } from '../models/device.model';

@Component({
  selector: 'app-select-machines',
  templateUrl: './select-machines.component.html',
  styleUrls: ['./select-machines.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  providers: [{
    provide: NG_VALUE_ACCESSOR,
    useExisting: SelectMachinesComponent,
    multi: true
  }]
})



export class SelectMachinesComponent implements OnInit, ControlValueAccessor {
  constructor(private data: DataService) { }

  currentMachine: Device;
  hidden = false;

  machines: Observable<Map<string, Device>>;
  private onChange: (_: any) => void = () => { };
  private onTouched: () => void = () => { };

  get value(): any {
    return this.currentMachine;
  }

  set value(v: any) {
    if (v !== this.currentMachine) {
      this.currentMachine = v;
      this.onChange(v);
    }
  }

  writeValue(val: Device): void {
    if (val != null && val!== undefined) {
      this.currentMachine = val;
    } else {
      this.onChange(this.currentMachine);
    }
  }
  registerOnChange(fn: any): void {
    this.onChange = fn;
  }
  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }
  setDisabledState(isDisabled: boolean): void {

  }
  onBlur() {
    this.onTouched();
  }



  ngOnInit() {
    this.data.getDevices();
    this.machines = this.data.devices$;
    this.machines.subscribe(data => {
      if (data.size <= 1) {
        this.hidden = true;
      }
      if (data.size > 0) {
        this.value = data.values().next().value;
        this.currentMachine = this.value;
      }
    });
  }

}
