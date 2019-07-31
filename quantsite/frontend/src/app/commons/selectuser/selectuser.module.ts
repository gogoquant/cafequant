import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';
import { SelectuserComponent } from './selectuser.component';
@NgModule({
  imports: [
    CommonModule,
    NgZorroAntdModule,
    FormsModule,
    ReactiveFormsModule
  ],
  exports: [
    SelectuserComponent
  ],
  providers: [{ provide: NZ_I18N, useValue: zh_CN }],
  declarations: [ SelectuserComponent ]
})
export class SelectuserModule { }
