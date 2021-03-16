import { Component, OnInit, ViewChild, ElementRef, Input, OnDestroy } from '@angular/core';
import { CycleSpec } from '../models/cyclespec.model';
import { Observable, Subject } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { UtilService, BroadcasterService } from '@menucha-de/shared';
import { ReportService } from './report.service';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-reports',
  templateUrl: './reports.component.html',
  styleUrls: ['./reports.component.scss']
})
export class ReportsComponent implements OnInit, OnDestroy {
  @ViewChild('uploader') fileDialogEl: ElementRef;

  @Input() group: string;
  items$: Observable<Map<string, CycleSpec>>;
  showContent = false;
  destroy$: Subject<boolean> = new Subject<boolean>();
  xx = '';
  constructor(
    private data: ReportService,
    private route: ActivatedRoute,
    private router: Router,
    private util: UtilService,
    private broadcaster: BroadcasterService) { }



  ngOnInit() {
    this.items$ = this.data.reports$;
    this.broadcaster.on('repadd').pipe(takeUntil(this.destroy$)).subscribe((xx: string) => {
      this.xx = xx;
    });
  }
  getshowv(id: string) {
    if (id === this.xx) {
      return true;
    }
    return false;
  }
  ngOnDestroy() {
    this.destroy$.next(true);
    this.destroy$.unsubscribe();
  }

  showReport(id: string) {
    this.router.navigate(['/report/' + id], { relativeTo: this.route });
  }

  importReport(files: FileList) {
    if (files && files.item.length > 0) {
      this.util.showSpinner();
      const reader = new FileReader();
      reader.readAsText(files.item(0));
      reader.onload = () => {
        try {
          const json = JSON.parse(reader.result as string);
          const xx = new CycleSpec(json);
          xx.id = '';
          xx.enabled = false;
          this.data.addReport(xx).subscribe(() => {
            this.util.hideSpinner();
          }, err => {
            let errMsg: string;
            const beg = err.error.indexOf('Unrecognized field');
            const end = err.error.indexOf('(') - 1;
            if (beg > 0 && end > beg) {
              errMsg = err.error.substring(beg, end);
            } else {
              errMsg = err.error;
            }
            this.util.showMessage('error', 'Failed to read file: ' + errMsg);
            this.util.hideSpinner();
          });
        } catch (ex) {
          console.log(ex)
          this.util.hideSpinner();
          this.util.showMessage('error', 'Wrong import file! ');
        }
      };
    }
    this.fileDialogEl.nativeElement.value = null;
  }
  refreshReports() {
    this.data.getReports();
  }

  onShow(value: boolean) {
    if (!value) {
      this.data.getReports();
    }
    this.showContent = !value;
  }
  trackByReportId(index: number, item) {
    return item.key;
  }
}
