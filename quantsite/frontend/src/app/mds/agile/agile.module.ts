import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AgileProductComponent } from './product/agile-product.component';
import { AgileProjectComponent } from './project/agile-project.component';
import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SelectuserModule } from 'src/app/commons/selectuser/selectuser.module';
@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    NgZorroAntdModule,
    ReactiveFormsModule,
    SelectuserModule
  ],
  declarations: [AgileProductComponent, AgileProjectComponent],
  exports: [
    // AgileProductComponent,
    // AgileProductComponent
  ]
})
export class AgileModule { }
