import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { TestComponent } from './test.component';
import { TmpComponent } from './tmp/tmp.component';
import { TtaComponent } from './tta/tta.component';

const testRoutes: Routes = [
    {
        path: '',
        component: TestComponent,
        children: [
            {
                path: 'tmp',
                component: TmpComponent,
                pathMatch: 'full',
                data: {
                    breadcrumb: '测试平台'
                }
            },
            {
                path: 'tta',
                component: TtaComponent,
                pathMatch: 'full',
                data: {
                    breadcrumb: '工具自动化'
                }
            },
            {
                path: '',
                redirectTo: 'tmp',
                pathMatch: 'full',
            }
        ],
        data: {
            breadcrumb: '测试管理'
        }
    }
];

@NgModule(
    {
        imports: [
            RouterModule.forChild(testRoutes)
        ],
        exports: [
            RouterModule
        ],
    }
)
export class TestRoutingModule {
}
