import { NgModule } from '@angular/core';
import { SelectMachinesComponent } from './select-machines.component';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
@NgModule({
  imports: [
    CommonModule,
    FormsModule
  ],
  declarations: [
    SelectMachinesComponent
  ],
  exports: [
    SelectMachinesComponent
  ]
})
export class SelectMachinesModule { }

