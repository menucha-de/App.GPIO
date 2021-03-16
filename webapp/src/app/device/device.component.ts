import { Component, OnInit, Input, OnDestroy, Inject } from '@angular/core';
import { DataService } from '../data.service';
import { Observable } from 'rxjs';


import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { Device } from '../models/device.model';
import { retryWhen, tap, delay } from 'rxjs/operators';
import { UtilService } from '@menucha-de/shared';

@Component({
  selector: 'app-device',
  templateUrl: './device.component.html',
  styleUrls: ['./device.component.scss']
})
export class DeviceComponent implements OnInit, OnDestroy {
  @Input() group: string;
  showContent = false;
  items$: Observable<Map<string, Device>>;
  myWebSocket: WebSocketSubject<any>;
  private configUrl: URL;
  reconnectdelay: 1000;

  constructor(
    private data: DataService,
    private util: UtilService,
    @Inject('WS_ENDPOINT') private wsEndpointTemplate: string,
  ) { }
  ngOnDestroy(): void {
    this.myWebSocket.unsubscribe();
    this.myWebSocket.complete();
  }

  ngOnInit(): void {
    this.items$ = this.data.devices$;
    this.configUrl = new URL(window.location.href);
    this.initWsConnection();
    /*let xx=1
    let timerId = setInterval(() =>{
      console.log('tick');
      this.data.mockvalue(1,xx)
      if (xx==1){
        xx=0;
      }
      else{
        xx=1;
      }
    }, 2000)*/
  }

  onShow(value: boolean) {
    if (!value) {
      this.data.getDevices();
    }
    this.showContent = !value;
  }
  createWebSocket(uri: string) {
    return new Observable(observer => {
      try {
        const subject = webSocket(uri);
        const subscription = subject.asObservable()
          .subscribe(data =>
            observer.next(data),
            error => observer.error(error),
            () => observer.complete());

        return () => {
          if (!subscription.closed) {
            subscription.unsubscribe();
          }
        };
      } catch (error) {
        observer.error(error);
      }
    });
  }

  private initWsConnection(): void {

    this.createWebSocket(this.websocketEndpoint)
      .pipe(
        retryWhen(errors =>
          errors.pipe(
            tap(err => {
              console.error(err);
            }),
            delay(1000)
          )
        )
      )
      .subscribe(msg => { this.data.updatePin(msg); }, err => console.error(err));
  }
  get websocketEndpoint() {
    const appName = this.configUrl.pathname;
    const secure = this.configUrl.protocol === 'https:';
    const protocol = secure ? 'wss' : 'ws';
    let result = this.wsEndpointTemplate.replace('{hostname}', this.configUrl.hostname);
    result = result.replace('{protocol}', protocol);
    result = result.replace('{container}', appName);
    return result;
  }
  changeLabel(item: Device) {
    if (item.label.trim() === '') {
      this.data.deleteLabel(item).subscribe(() => null
        , err => {
          this.util.showMessage('error', err.error);

        });
    } else {
      this.data.updateLabel(item).subscribe(() => null
        , err => {
          this.util.showMessage('error', err.error);

        });
    }
  }
}
