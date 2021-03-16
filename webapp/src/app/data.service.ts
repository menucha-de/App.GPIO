import { Injectable } from '@angular/core';
import { Observable, BehaviorSubject, throwError } from 'rxjs';
import { CommonModule } from '@angular/common';
import { UtilService } from '@menucha-de/shared';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';

import { Device } from './models/device.model';
import { Field } from './models/field.model';

@Injectable({
  providedIn: 'root'
})

@Injectable({
  providedIn: CommonModule
})
export class DataService {

  private readonly baseUrl = 'rest/gpio/';
  devices$: Observable<Map<string, Device>>;
  private _devices: BehaviorSubject<Map<string, Device>>;
  private dataStore: {
    devices: Map<string, Device>;
  };

  constructor(protected http: HttpClient, private util: UtilService) {
    this.dataStore = {
      devices:
        new Map<string, Device>()
    };
    this._devices = new BehaviorSubject(new Map<string, Device>());
    this.devices$ = this._devices.asObservable();

  }
  getDevices() {
    return this.http.get<Map<string, Device>>(`${this.baseUrl}devices`).subscribe((data: Map<string, Device>) => {
      this.dataStore.devices = new Map<string, Device>();
      Object.entries(data).forEach(([key, value]) => {
        if (value.label == null) { value.label = value.name; }
        const dev = new Device(value);

        dev.setProperties(value.properties);
        dev.setFields(value.fields);
        this.dataStore.devices.set(key, dev);

      });
      this._devices.next(Object.assign({}, this.dataStore).devices);
    }, err => this.handleError(err));
  }


  updatePin(pin) {
    const field = new Field(pin.field);
    if (this.dataStore.devices.has(pin.deviceId)) {
      if (field.properties.size > 0) {

        if (!this.dataStore.devices.get(pin.deviceId).fields.get(field.id).equals(field)) {
          this.dataStore.devices.get(pin.deviceId).fields.set(field.id, field);
          this._devices.next(Object.assign({}, this.dataStore).devices);

        }

      } else { // only update value
        const xx = this.dataStore.devices.get(pin.deviceId).fields.get(field.id);
        if (xx.value !== field.value) {
          this.dataStore.devices.get(pin.deviceId).fields.get(field.id).value = field.value;
        }
      }
    }
  }
  setProperty(device: string, field: string, name: string, data: string) {
    return this.http.put(`${this.baseUrl}devices/${device}/fields/${field}/properties/${name}`, data,
      { headers: { 'content-type': 'text/plain' } });
  }
  setLabel(device: string, field: string, data: string) {
    return this.http.put(`${this.baseUrl}devices/${device}/fields/${field}/label`, data,
      { headers: { 'content-type': 'text/plain' } });
  }
  deleteFieldLabel(device: string, field: string) {
    return this.http.delete(`${this.baseUrl}devices/${device}/fields/${field}/label`);
  }
  updateLabel(device: Device) {
    return this.http.put(`${this.baseUrl}devices/${device.id}/label`, device.label,
      { headers: { 'content-type': 'text/plain' } });
  }
  deleteLabel(device: Device) {
    return this.http.delete(`${this.baseUrl}devices/${device.id}/label`);
  }
  setStatus(device: string, field: string, data: boolean) {
    let xx = 'LOW';
    if (data) {
      xx = 'HIGH';
    }
    return this.http.put(`${this.baseUrl}devices/${device}/fields/${field}/value`, xx,
      { headers: { 'content-type': 'text/plain' } });
  }
  getDevice(device: string) {
    return this.http.get<Device>(`${this.baseUrl}devices/${device}`);
  }

  private handleError(error: HttpErrorResponse) {
    this.util.showMessage('error', error.error);
    return throwError(error);
  }
  mockvalue(pin, data) {
    this.http.get(`${this.baseUrl}setpinvalue/${pin}/value/${data}`).subscribe();
  }
}
