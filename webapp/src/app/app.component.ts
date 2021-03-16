import { Component, Inject } from '@angular/core';



@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'GPIO';
  constructor(
    @Inject('SHOW_FRAME') private _showFrame: boolean
  ) {}

  get showFrame() {
    return this._showFrame;
  }

}
