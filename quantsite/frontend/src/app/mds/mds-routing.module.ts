import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {MdsComponent} from './mds.component';
import { AgileProductComponent } from './agile/product/agile-product.component';
import { AgileProjectComponent } from './agile/project/agile-project.component';

const routes: Routes = [
  {
    path: 'mds',
    component : MdsComponent,
    children: [
      {
        path: 'projects',
        component: AgileProjectComponent,
        pathMatch: 'full'
      },
      {
        path: 'products',
        component: AgileProductComponent,
        pathMatch: 'full'
      },
      {
        path: 'test',
        loadChildren: './test/test.module#TestModule',
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MdsRoutingModule { }
