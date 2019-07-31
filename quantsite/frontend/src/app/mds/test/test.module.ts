import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { TmpModule } from './tmp/tmp.module';
import { TtaModule } from './tta/tta.module';
import { TestComponent } from './test.component';
import { TestRoutingModule } from './test-routing.module';
import { BreadcrumbComponent } from './breadcrumb/breadcrumb.component';
import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';

@NgModule({
  imports: [
    RouterModule,
    CommonModule,
    NgZorroAntdModule,
    TmpModule,
    TtaModule,
    TestRoutingModule,
  ],
  providers: [{ provide: NZ_I18N, useValue: zh_CN }],
  declarations: [TestComponent, BreadcrumbComponent]
})
export class TestModule { }
