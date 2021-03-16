import { Component, forwardRef, ElementRef, ViewChild, Input } from '@angular/core';
import { NG_VALUE_ACCESSOR, ControlValueAccessor } from '@angular/forms';


const noop = () => {
};

export const SWITCH_VALUE_ACCESSOR: any = {
  provide: NG_VALUE_ACCESSOR,
  useExisting: forwardRef(() => CustomSwitchComponent),
  multi: true
};

let identifier = 0;

@Component({
  selector: 'app-customswitch',
  templateUrl: './customswitch.component.html',
  styleUrls: ['./customswitch.component.scss'],
  providers: [SWITCH_VALUE_ACCESSOR]
})
export class CustomSwitchComponent implements ControlValueAccessor {

  @ViewChild('check') check: ElementRef;
  @Input() values: string[];
  disabled: boolean;
  private _value: string;
  private _val: boolean;
  private state = 1;
  identifier = `sw-${identifier++}`;

  private onTouchedCallback: () => void = noop;
  private onChangeCallback: (_: any) => void = noop;

  get val(): boolean {
    return this._val;
  }

  set val(value: boolean) {
    const el = this.check.nativeElement;
    let xx: string;
    switch (this.state) {
      case 0:
        el.indeterminate = true;
        this.state++;
        xx = this.values[1];
        break;
      case 1:
        el.indeterminate = false;
        this.state++;
        el.checked = true;
        xx = this.values[2];
        break;
      case 2:
        el.checked = false;
        this.state = 0;
        xx = this.values[0];
        break;
    }
    if (xx !== this._value) {
      this._value = xx;
      this.onChangeCallback(xx);
    }
  }

  set value(value: string) {
    if (value !== this._value) {
      this._value = value;
      this.onChangeCallback(value);
    }
  }

  onBlur() {
    this.onTouchedCallback();
  }

  writeValue(value: any): void {//must be run twice for viewinit
    this._value = value;

    this.state = this.values.indexOf(value);

    if (this.check != null) {
      if (this.state === 1) {
        this.check.nativeElement.indeterminate = true;
        this.check.nativeElement.checked = false;
      } else {
        this.check.nativeElement.indeterminate = false;
        if (this.state === 0) {
          this._val = false;
        } else {
          this._val = true;
        }
      }
    }
  }

  registerOnChange(fn: any): void {
    this.onChangeCallback = fn;
  }
  registerOnTouched(fn: any): void {
    this.onTouchedCallback = fn;
  }
  setDisabledState?(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }

  constructor() { }
  stateChange(xx: any) {
    const el = xx.currentTarget;
    el.indeterminate = false;
    switch (this.state) {
      case 0:
        el.indeterminate = true;
        this.state++;

        break;
      case 1:
        el.indeterminate = false;
        this.state++;
        el.checked = true;
        break;
      case 2:
        el.checked = false;
        this.state = 0;
        break;
    }
  }
}
