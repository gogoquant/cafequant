import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { MdsRoutingModule } from './mds/mds-routing.module';
import { AccountComponent } from './account/account.component';
import { AccountgroupComponent } from './account/group/account-group.component';
import { LoggingComponent } from './logging/logging.component';
import { DashboardComponent } from './system/dashboard/dashboard.component';
import { MceComponent } from './mce/mce.component';
import { ClusterComponent } from './mce/cluster/cluster.component';
import { MseComponent } from './mse/mse.component';
import { TenantComponent } from './account/tenant/tenant.component';

const routes: Routes = [
  { path: 'mce', component: MceComponent, pathMatch: 'full' },
  { path: 'cluster', component: ClusterComponent, pathMatch: 'full' },
  { path: 'users', component: AccountComponent, pathMatch: 'full' },
  { path: 'usersgroup', component: AccountgroupComponent, pathMatch: 'full' },
  { path: 'tenant', component: TenantComponent, pathMatch: 'full' },
  { path: 'mse', loadChildren: './mse/mse.module#MseModule' },
  { path: 'mds', loadChildren: './mds/mds.module#MdsModule' },
  // { path: 'mse', component: MseComponent, pathMatch: 'full', data: {breadcrumb: '配置文件中心'}},
  { path: 'logging', component: LoggingComponent, pathMatch: 'full' },
  { path: 'dashboard', component: DashboardComponent, pathMatch: 'full' },
  { path: '', redirectTo: '/users', pathMatch: 'full' },
  { path: '**', redirectTo: '/', pathMatch: 'full' }
];
@NgModule({
  imports: [
    RouterModule.forRoot(routes,
      {
        enableTracing: true,
      })
  ],
  exports: [RouterModule],
  declarations: []
})
export class AppRoutingModule { }
