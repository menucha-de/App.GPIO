import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReportComponent } from './report/report.component';
import { ReportsComponent } from './reports.component';
import { ReportDetailComponent } from './report-detail/report-detail.component';
import { SubscriptionsComponent } from './subscriptions/subscriptions.component';
import { MicaAppComponentsModule } from '@peramic/shared';
import { MicaControlsModule } from '@peramic/controls';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SelectMachinesModule } from '../select-machines/select-machines.module';




@NgModule({
  declarations: [
    ReportComponent,
    ReportsComponent,
    ReportDetailComponent,
    SubscriptionsComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MicaControlsModule,
    MicaAppComponentsModule,
    SelectMachinesModule
  ],
  exports: [
    ReportComponent,
    ReportsComponent,
    ReportDetailComponent,
    SubscriptionsComponent
  ]
})
export class ReportsModule { }
