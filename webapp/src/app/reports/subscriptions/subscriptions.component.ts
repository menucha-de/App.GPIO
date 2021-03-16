import { Component, OnInit, Input, OnDestroy } from '@angular/core';

import { Router } from '@angular/router';
import { BroadcasterService, Subscriptor, UtilService, Subscriber, TransportService } from '@menucha-de/shared';
import { ServiceState } from '@menucha-de/controls';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ReportService } from '../report.service';

@Component({
  selector: 'app-subscriptions',
  templateUrl: './subscriptions.component.html',
  styleUrls: ['./subscriptions.component.scss']
})
export class SubscriptionsComponent implements OnInit, OnDestroy {
  @Input() cycleId: string;

  @Input() disabled: boolean;
  destroy$: Subject<boolean> = new Subject<boolean>();
  serviceSubscriptorState: Map<string, ServiceState> = new Map<string, ServiceState>();

  constructor(
    private data: ReportService,
    private broadcaster: BroadcasterService,
    private router: Router,
    private util: UtilService,
    private transport: TransportService) { }
  subscriptions: Map<string, Subscriptor>;
  subscriptors: Map<string, string> = new Map();

  ngOnInit() {
    this.getSubscriptions();
    this.broadcaster.on('subscriptorsChanged').pipe(takeUntil(this.destroy$)).subscribe(() => {
      this.getSubscriptions();
    }, err => {
      this.util.showMessage('error', err.error);
    });
    this.broadcaster.on('subsChanged').pipe(takeUntil(this.destroy$)).subscribe((xx) => {
      if (xx === this.cycleId) { this.getSubscriptions(); }

    }, err => {
      this.util.showMessage('error', err.error);
    });
  }
  ngOnDestroy() {
    this.destroy$.next(true);
    this.destroy$.unsubscribe();
  }
  getSubscriptions() {
    this.data.getSubscriptions(this.cycleId).pipe(takeUntil(this.destroy$)).subscribe((data) => {
      if (data != null && data !== undefined) {
        this.subscriptions = new Map(Object.entries(data));
        this.subscriptions.forEach(element => {
          this.transport.getSubscriber(element.subscriberId).subscribe(xx => {
            this.subscriptors.set(xx.id, xx.uri);
          });
        });
      }
    }, this.errorHandler);
  }
  addSubscriptor() {
    this.router.navigate(['/subscriptions/new',
      { suffix: this.cycleId + '/subscriptions' }]);
  }
  deleteSubscriptor(id: string) {
    this.data.deleteSubscription(this.cycleId, id).pipe(takeUntil(this.destroy$)).subscribe(() => {
      this.subscriptions.delete(id);
    }, this.errorHandler);
  }
  onChangeSubscriptor(subscriptorId: string) {
    this.router.navigate(['/subscriptions/' + subscriptorId,
    { suffix: this.cycleId + '/subscriptions' }]);
  }
  enableSubscriptor(subscriptor: Subscriptor) {
    if (this.isSubscriptorPending(subscriptor)) {
      return;
    }
    this.serviceSubscriptorState.set(subscriptor.id, ServiceState.Pending);
    subscriptor.enable = !subscriptor.enable;
    this.data.updateSubscription(this.cycleId, subscriptor.id, subscriptor).pipe(takeUntil(this.destroy$)).subscribe(() => {
      this.serviceSubscriptorState.set(subscriptor.id, this.getSubscriptorState(subscriptor));
    }, (err) => {
      this.util.showMessage('error', err.error);
      subscriptor.enable = !subscriptor.enable;
      this.serviceSubscriptorState.set(subscriptor.id, this.getSubscriptorState(subscriptor));
    });
  }
  getSubscriptorState(subscriptor: Subscriptor) {
    return (subscriptor && (subscriptor.enable === true)) ? ServiceState.Started : ServiceState.Stopped;
  }

  isSubscriptorPending(subscriptor: Subscriptor) {
    const serviceState = this.serviceSubscriptorState.get(subscriptor.id);
    return serviceState === ServiceState.Pending;
  }

  private errorHandler = (err: { error: string; }) => {
    this.util.showMessage('error', err.error);
  }
}
