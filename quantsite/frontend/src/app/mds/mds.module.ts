import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AgileModule } from './agile/agile.module';
import { CicdModule } from './cicd/cicd.module';
import { TestModule } from './test/test.module';
import { MdsComponent } from './mds.component';
import {MdsRoutingModule} from './mds-routing.module';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';

@NgModule({
  imports: [
    CommonModule,
    NgZorroAntdModule,
    AgileModule,
    CicdModule,
    TestModule,
    FormsModule,
    ReactiveFormsModule,
    MdsRoutingModule,
  ],
  providers: [{ provide: NZ_I18N, useValue: zh_CN }],
  declarations: [MdsComponent]
})
export class MdsModule {

}
