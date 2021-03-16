import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, BehaviorSubject, of, throwError } from 'rxjs';

import { tap, map } from 'rxjs/operators';
import { Subscriptor, UtilService } from '@menucha-de/shared';
import { CycleSpec } from '../models/cyclespec.model';
import { Device } from '../models/device.model';

@Injectable({
    providedIn: 'root'
})
export class ReportService {
    private static readonly BASE_PATH = 'rest/gpio/';
    reports$: Observable<Map<string, CycleSpec>>;
    private _reports: BehaviorSubject<Map<string, CycleSpec>>;
    private reportStore: {
        reports: Map<string, CycleSpec>
    };
    constructor(private http: HttpClient, private util: UtilService) {
        this.reportStore = {
            reports: new Map<string, CycleSpec>()
        };
        this._reports = new BehaviorSubject(new Map<string, CycleSpec>());
        this.reports$ = this._reports.asObservable();
    }
    getReports() {
        this.http.get<Map<string, CycleSpec>>(ReportService.BASE_PATH + `reports`).subscribe((data: Map<string, CycleSpec>) => {
            if (data !== undefined && data != null) {
                Object.entries(data).forEach(([key, value]) => {
                    const rep = new CycleSpec(value);
                    this.reportStore.reports.set(key, rep);
                });
            }
            this._reports.next(Object.assign({}, this.reportStore).reports);
        }, err => this.handleError(err));
    }
    getReport(id: string) {
        if (id === 'new') {
            return of(new CycleSpec());
        }
        return this.http.get<CycleSpec>(ReportService.BASE_PATH + `reports/${id}`, { observe: 'response' }).pipe(map(response => {

            if (response.status === 204) {
                return new CycleSpec();
            } else {
                return new CycleSpec(response.body);

            }
        }));

    }
    addReport(item: CycleSpec) {
        return this.http.post(ReportService.BASE_PATH + `reports`, item.toObject(), { responseType: 'text' })
            .pipe(tap((data: string) => {
                item.id = data;
                this.reportStore.reports.set(data, item);
                this._reports.next(Object.assign({}, this.reportStore).reports);
            }));
    }
    setReport(id: string, report: CycleSpec) {
        return this.http.put(ReportService.BASE_PATH + `reports/${id}`, report.toObject())
            .pipe(tap(() => {
                this.reportStore.reports.set(id, report);
                this._reports.next(Object.assign({}, this.reportStore).reports);
            }));
    }
    deleteReport(id: string) {
        return this.http.delete(ReportService.BASE_PATH + `reports/${id}`)
            .pipe(tap(() => {
                this.reportStore.reports.delete(id);
                this._reports.next(Object.assign({}, this.reportStore).reports);
            }));
    }
    getSubscriptions(id: string) {
        return this.http.get<Map<string, Subscriptor>>(ReportService.BASE_PATH + `reports/${id}/subscriptions`);
    }
    deleteSubscription(reportId: string, id: string) {
        reportId = encodeURI(reportId);
        return this.http.delete(ReportService.BASE_PATH + `reports/${reportId}/subscriptions/${id}`);
    }
    updateSubscription(id: string, subId: string, subscriptor: Subscriptor) {

        return this.http.put(ReportService.BASE_PATH + `reports/${id}/subscriptions/${subId}`, subscriptor);
    }
    private handleError(error: HttpErrorResponse) {
        this.util.showMessage('error', error.error);
        return throwError(error);
    }

}
