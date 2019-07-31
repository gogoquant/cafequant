import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { MseComponent } from './mse.component';
import { AppInfoListComponent } from './app-info-list/app-info-list.component';
import { AppInfoDetailComponent } from './app-info-detail/app-info-detail.component';

const mseRoutes: Routes = [
    {
        path: '',
        component: MseComponent,
        data: {
            breadcrumb: '应用配置中心'
        },
        children: [
            {
                path: 'apps',
                component: AppInfoListComponent,
                data: {
                    breadcrumb: '应用配置列表'
                }
            },
            {
                path: 'app/:id',
                component: AppInfoDetailComponent,
                data: {
                    breadcrumb: '应用配置详情'
                }
            },
            {
                path: '',
                redirectTo: 'apps',
                pathMatch: 'full'
            }
        ]
    }
];

@NgModule(
    {
        imports: [
            RouterModule.forChild(mseRoutes)
        ],
        exports: [
            RouterModule
        ],
    }
)
export class MseRoutingModule {
}
