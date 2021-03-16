import { Component, OnInit } from '@angular/core';

import { Observable, of } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';

import { switchMap, tap, catchError } from 'rxjs/operators';
import { CloseAction, UtilService, BroadcasterService } from '@peramic/shared';


import { CycleSpec } from '../../models/cyclespec.model';
import { ReportService } from '../report.service';
import { PinUI } from 'src/app/models/pinui.model';
import { Field } from 'src/app/models/field.model';
import { Device } from 'src/app/models/device.model';

@Component({
  selector: 'app-report-detail',
  templateUrl: './report-detail.component.html',
  styleUrls: ['./report-detail.component.scss']
})
export class ReportDetailComponent implements OnInit {

  fieldSubscriptions: Field[] = [];
  oldFieldSubscriptions: string[] = [];
  report$: Observable<CycleSpec>;
  device: Device;
  id: string;

  constructor(
    private route: ActivatedRoute,
    private util: UtilService,
    private router: Router,
    private data: ReportService,
    private broadcaster: BroadcasterService
  ) { }

  ngOnInit() {
    this.fieldSubscriptions = [];
    this.report$ = this.route.paramMap.pipe(switchMap(paramMap => {
      this.id = paramMap.get('id');
      return this.data.getReport(this.id);
    }), tap(report => {
      report.fieldSubscriptions.forEach((fieldIds) => {
        fieldIds.forEach(fieldId => {
          const aa = new Field();
          aa.id = fieldId;
          this.fieldSubscriptions.push(aa);
          this.oldFieldSubscriptions.push(fieldId);
        });

      });
    }), catchError((err) => {
      this.util.showMessage('error', err.error);
      this.router.navigate(['']);
      return of(null);
    }));
  }

  reportEnable(report: CycleSpec) {
    if (report.enabled == null) {
      report.enabled = true;
    } else {
      report.enabled = !report.enabled;
    }
  }

  onClose(closeAction: CloseAction, item: CycleSpec) {
    if (closeAction === CloseAction.OK) {
      if (!item.name || item.name.trim().length === 0) {
        this.util.showMessage('error', 'Name must not be empty');
        return;
      } else if ((item.duration == null) || (item.duration < 0)) {
        this.util.showMessage('error', 'Invalid duration value');
        return;
      } else if ((item.repeatPeriod == null) || (item.repeatPeriod < 0)) {
        this.util.showMessage('error', 'Invalid repeat period value');
        return;
      }
      if ((item.duration == null || item.duration <= 0) &&
        (item.repeatPeriod == null || item.repeatPeriod <= 0) &&
        !item.whenDataAvailable) {
        this.util.showMessage('error', 'No stop condition specifed');
        return;
      }
      item.fieldSubscriptions = new Map<string, Set<string>>();
      const fields: Set<string> = new Set();
      this.fieldSubscriptions.forEach((deviceField: Field) => {
        if ((deviceField.id !== '0') && (deviceField.id != null)) {
          fields.add(deviceField.id);
        }
      });
      if (fields.size > 0) {
        item.fieldSubscriptions.set(this.device.id, fields);
      }
      const saveReport$ = this.id !== 'new'
        ? this.data.setReport(this.id, item)
        : this.data.addReport(item);
      saveReport$.subscribe(() => {
        if (this.id === 'new') {
          this.broadcaster.broadcast('repadd', item.id);
        }
        this.util.hideSpinner();
        this.router.navigate(['']);
      }, err => {
        this.util.hideSpinner();
        this.util.showMessage('error', err.error);
      });

    } else {
      this.router.navigate(['']);
    }
  }
  remove(deviceField: Field) {
    const indexOf = this.fieldSubscriptions.findIndex(df => df.id === deviceField.id);

    if (indexOf > -1) {
      this.fieldSubscriptions.splice(indexOf, 1);
      this.oldFieldSubscriptions.splice(indexOf, 1);
    }
  }

  addField() {
    this.fieldSubscriptions.push(new PinUI(null));
  }

  getDeviceFields(dev: Device): Field[] {
    const result: Field[] = [];
    if (dev != null) {
      if (dev.fields != null) {
        dev.fields.forEach((field) => {
          const pin = new PinUI(field);
          if (pin.direction === 'INPUT') {
            result.push(field);
          }
        });
      }
    }
    return result;
  }

  onEditField(pin: Field, items: Field[]) {

    if (pin.id === '0') {
      return;
    }
    if (pin.id === 'ADD ALL FIELDS') {
      this.fieldSubscriptions = [];
      this.oldFieldSubscriptions = [];
      items.forEach((el) => {
        this.fieldSubscriptions.push(el);
        this.oldFieldSubscriptions.push(el.id);
      });


    } else {

      const indexOf = this.oldFieldSubscriptions.findIndex(fi => fi === pin.id);
      if (indexOf >= 0) {
        this.fieldSubscriptions.splice(indexOf, 1);
      }
      this.oldFieldSubscriptions = this.fieldSubscriptions.map(t => t.id);
    }
  }
}
