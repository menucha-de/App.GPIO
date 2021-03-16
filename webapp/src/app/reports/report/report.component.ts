import { Component, OnInit, Input } from '@angular/core';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { BroadcasterService, UtilService } from '@peramic/shared';
import { CycleSpec } from '../../models/cyclespec.model';
import { Router } from '@angular/router';
import { ReportService } from '../report.service';
import { ServiceState } from '@peramic/controls';

@Component({
  selector: 'app-report',
  templateUrl: './report.component.html',
  styleUrls: ['./report.component.scss'],
  animations: [
    trigger('showContent', [
      state('true', style({ height: '*' })),
      state('false', style({ height: 0 })),
      transition('true => false', animate('400ms ease-out')),
      transition('false => true', animate('400ms ease-in'))
    ]),
  ]
})
export class ReportComponent implements OnInit {
  @Input() report: CycleSpec;
  @Input() id: string;
  @Input() showv: boolean;
  showSubscriptors = false;
  reportState: ServiceState = ServiceState.Stopped;

  constructor(
    private router: Router,
    private broadcaster: BroadcasterService,
    private data: ReportService,
    private util: UtilService) { }

  ngOnInit() {
    this.reportState = this.getReportState(this.report);
    this.showSubscriptors = this.showv;
  }

  getReportState(cycleSpec: CycleSpec) {
    return (cycleSpec && (cycleSpec.enabled === true)) ? ServiceState.Started : ServiceState.Stopped;
  }
  isReportPending() {
    return this.reportState === ServiceState.Pending;
  }

  showReport() {
    this.router.navigate(['/report/', this.id]);
  }

  exportReport() {
    const res = new Blob([JSON.stringify(this.report.toObject())]);
    const filename = 'report_' + this.report.name + '.json';
    if (navigator.appVersion.toString().indexOf('.NET') > 0) {
      window.navigator.msSaveBlob(res, filename);
    } else {
      const url = window.URL.createObjectURL(res);
      const a = document.createElement('a');
      document.body.appendChild(a);
      a.setAttribute('style', 'display: none');
      a.href = url;
      a.download = filename;
      a.click();
      window.URL.revokeObjectURL(url);
      a.remove();
    }
  }

  undefineReport() {
    this.data.deleteReport(this.id).subscribe(() => null, this.errorHandler);

  }
  setReportState() {
    if (this.isReportPending()) {
      return;
    }
    this.reportState = ServiceState.Pending;
    this.report.enabled = !this.report.enabled;
    this.data.setReport(this.id, this.report).subscribe(() => {
      this.reportState = this.getReportState(this.report);

      if (this.report.enabled) {
        this.broadcaster.broadcast('subsChanged', this.report.name);
      }
    },
      err => {
        this.util.showMessage('error', err.error);
        this.report.enabled = !this.report.enabled;
        this.reportState = this.getReportState(this.report);
      }
    );
  }
  onDefaultAction() {
    this.showSubscriptors = !this.showSubscriptors;
  }

  private errorHandler = (err: { error: string; }) => {
    this.util.showMessage('error', err.error);
  }
}
