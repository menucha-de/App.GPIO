import { Component, OnInit, Input } from '@angular/core';

import { PinUI } from '../models/pinui.model';
import { DataService } from '../data.service';
import { UtilService } from '@menucha-de/shared';
import { Field } from '../models/field.model';

@Component({
  selector: 'app-pin',
  templateUrl: './pin.component.html',
  styleUrls: ['./pin.component.scss']
})
export class PinComponent implements OnInit {
  @Input() field: Field;
  @Input() id: string;
  @Input() state: boolean;
  pin: PinUI;
  defaultvalues = ['LOW', null, 'HIGH'];
  oldlabel = '';
  constructor(private data: DataService, private util: UtilService) { }
  ngOnInit(): void {
    this.pin = new PinUI(this.field);
    this.oldlabel = this.pin.label;
  }

  statusChange() {
    this.data.setStatus(this.id, this.field.id, this.state).subscribe(() => null
      , err => {
        this.util.showMessage('error', err.error);
        this.state = !this.state;
      });
  }
  directionChange() {
    this.data.setProperty(this.id, this.field.id, 'direction', this.pin.direction).subscribe(() => {
    }
      , err => {
        this.util.showMessage('error', err.error);
        if (this.pin.direction === 'INPUT') {
          this.pin.direction = 'OUTPUT';
        } else {
          this.pin.direction = 'INPUT';
        }
      });
  }
  initialStatusChange() {
    let data = this.pin.initialState;
    if (this.pin.initialState == null) {
      data = 'UNKNOWN';
    }
    this.data.setProperty(this.id, this.field.id, 'initialState', data).subscribe(() => null
      , err => {
        this.util.showMessage('error', err.error);
      });
  }
  changeLabel() {
    if (this.pin.label.trim() === '') {
      this.data.deleteFieldLabel(this.id, this.field.id).subscribe(() => {
        this.oldlabel = this.pin.label;
      }
        , err => {
          this.pin.label = this.oldlabel;
          this.util.showMessage('error', err.error);
        });
    } else {
      this.data.setLabel(this.id, this.field.id, this.pin.label).subscribe(() => {
        this.oldlabel = this.pin.label;
      }
        , err => {
          this.pin.label = this.oldlabel;
          this.util.showMessage('error', err.error);
        });
    }
  }
}
