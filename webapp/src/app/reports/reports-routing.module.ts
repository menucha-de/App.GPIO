import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ReportDetailComponent } from './report-detail/report-detail.component';


const routes: Routes = [{
    path: 'report',
    children: [
        {
            path: ':id',
            component: ReportDetailComponent
        }
    ]

}];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class ReportsRoutingModule { }
