import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { transportRoutes, SubscriptionDialogComponent } from '@menucha-de/shared';



const routes: Routes = [{ path: '', redirectTo: 'main', pathMatch: 'full' },
{
  path: 'subscriptions/:id', component: SubscriptionDialogComponent
},
 ...transportRoutes
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
