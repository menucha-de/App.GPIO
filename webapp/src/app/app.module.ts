import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MicaAppBaseModule, TransportModule, TransportConfig, MicaAppComponentsModule } from '@peramic/shared';
import { MicaControlsModule } from '@peramic/controls';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';

import { MainComponent } from './main/main.component';
import { FormsModule } from '@angular/forms';
import { PinComponent } from './pin/pin.component';
import { DeviceComponent } from './device/device.component';
import { CustomSwitchComponent } from './customswitch/customswitch.component';

import { ReportsRoutingModule } from './reports/reports-routing.module';
import { ReportsModule } from './reports/reports.module';
import { SelectMachinesModule } from './select-machines/select-machines.module';


const transportConfig: TransportConfig = {
  subscribers: {
    restBaseUrl: '',
    routeParentUrl: ''
  },
  subscriptions: {
    restBaseUrl: 'gpio/reports/',
    routeParentUrl: ''
  }
};


@NgModule({
  declarations: [
    AppComponent,
    MainComponent,
    PinComponent,
    DeviceComponent,
    CustomSwitchComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    MicaAppBaseModule,
    MicaAppComponentsModule,
    MicaControlsModule,
    HttpClientModule,
    BrowserAnimationsModule,
    FormsModule,
    SelectMachinesModule,
    TransportModule.forRoot(transportConfig),
    ReportsModule,
    ReportsRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
